package config

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"context"
	"os/exec"

	"github.com/b4b4r07/go-finder"
)

type (
	CdxSource struct {
		Name       string `yaml:"name"`
		Alias      string `yaml:"alias"`
		Command    string `yaml:"command"`
		SkipColumn int    `yaml:"skipColumn"`
	}

	Collection []CdxSource
)

func (s CdxSource) run(ctx context.Context, listener chan<- finder.Item) error {
	ch := make(chan error)
	go func() {
		cmd := exec.Command(DefaultShell(), "-c", s.Command)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- err
			return
		}

		err = cmd.Start()
		if err != nil {
			ch <- err
			return
		}

		scan := bufio.NewScanner(stdout)
		for scan.Scan() {
			text := scan.Text()
			listener <- finder.Item{
				Key:   fmt.Sprintf("[%s]\t%s", s.Name, text),
				Value: strings.Fields(text)[s.SkipColumn:],
			}
		}

		ch <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		return errors.New("canceled by user")
	case err := <-ch:
		return err
	}
}

// NewCollection create Collection struct from config and cli option
// params:
// 	- b: /path/to/bookmark
//  - h: /path/to/history
//  - cfg: defined Collection by config file
//  - req: requests from cli
// returns:
//  - c: result collection
//  - err: error
func NewCollection(h, b string, cfg Collection, req string) (c Collection, err error) {

	m := map[string]CdxSource{
		"h": {Name: "history", Alias: "h", SkipColumn: 0, Command: fmt.Sprintf("cat %s", h)},
		"b": {Name: "bookmark", Alias: "b", SkipColumn: 0, Command: fmt.Sprintf("cat %s", b)},
	}

	for _, v := range cfg {
		m[v.Alias] = v
	}

	c = Collection{}
	for _, v := range req {
		if s, ok := m[string(v)]; ok {
			c = append(c, s)
		} else {
			return nil, errors.Errorf("cannot find cdxsource: %s", s.Alias)
		}
	}

	return
}

// generateFinderItem
// param:
// 	- ctx: context
//  - fromCli: paths from cli
// returns:
// 	- items: finder items
// 	- err: error
func (c Collection) generateFinderItem(ctx context.Context, fromCli []string) (items finder.Items, err error) {
	// receiver for cdxsource.run
	listener := make(chan finder.Item, 20)
	// error channel
	ch := make(chan error, 1)
	defer close(ch)

	var wg sync.WaitGroup
	wg.Add(len(c))

	// start source commands
	for _, v := range c {
		go func(s CdxSource) {
			if err := s.run(ctx, listener); err != nil {
				ch <- err
			}
			wg.Done()
		}(v)
	}

	go func() {
		wg.Wait()
		close(listener)
	}()

	items = finder.Items{}
	for _, v := range fromCli {
		items.Add(v, []string{v})
	}

	for {
		select {
		case err = <-ch:
			return
		case item, ok := <-listener:
			if ok {
				items.Add(item.Key, item.Value)
			} else {
				return
			}
		}
	}
}

// Select is select destination path
// params:
//	- ctx: context
//  - ff: fuzzy finder setting
// returns:
//  - path: destination path
//  - err: error
func (c Collection) Select(ctx context.Context, ff FuzzyFinder, paths []string) (path string, err error) {
	f, err := finder.New(append([]string{ff.Command}, ff.Options...)...)
	if err != nil {
		return
	}

	source, err := c.generateFinderItem(ctx, paths)
	if err != nil {
		return
	}

	res, err := f.Select(source)
	if err != nil {
		return
	}
	return strings.Join(res[0].([]string), " "), nil
}
