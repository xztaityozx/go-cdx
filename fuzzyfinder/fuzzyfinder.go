package fuzzyfinder

import (
	"context"
	"github.com/b4b4r07/go-finder"
	"github.com/b4b4r07/go-finder/source"
	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/go-cdx/customsource"
	"golang.org/x/xerrors"
	"sync"
)

type (
	FuzzyFinder struct {
		Path    string
		Options []string
	}
)

func (ff FuzzyFinder) Start(ctx context.Context, cs []customsource.CustomSource) ([]string, error) {
	f, err := finder.New(append([]string{ff.Path}, ff.Options...)...)
	if err != nil {
		return nil, xerrors.Errorf("failed build fuzzy finder: %w", err)
	}

	listener := make(chan string, 10)

	// start source commands
	var wg sync.WaitGroup
	for _, v := range cs {
		wg.Add(1)
		go func(v customsource.CustomSource) {
			defer wg.Done()
			res, err := v.Start()
			if err != nil {
				logrus.WithError(err).Errorf("failed CustomSource: %s", v.Name)
				return
			}

			for _, s := range res {
				if len(s) == 0 {
					continue
				}
				listener <- s
			}
		}(v)
	}

	var box []string
	ch := make(chan struct{})
	go func() {
		for v := range listener {
			box = append(box, v)
		}
		ch <- struct{}{}
	}()

	logrus.Info("waiting")

	wg.Wait()
	close(listener)

	select {
	case <-ctx.Done():
		return nil, xerrors.New("canceled")
	case <-ch:
	}

	f.Read(source.Slice(box))
	items, err := f.Run()

	return items, err
}
