package customsource

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"text/tabwriter"

	"github.com/b4b4r07/go-finder"
	"golang.org/x/xerrors"
)

type SourceCollection []CustomSource

// Print はCustomSourceを一覧表示する
// returns:
//  - error:
func (sc SourceCollection) Print() error {
	w := tabwriter.NewWriter(os.Stderr, 0, 8, 0, '\t', 0)
	_, err := fmt.Fprint(w, "Name\tAlias\tBeginColumn\tCommands")
	if err != nil {
		return err
	}
	for _, v := range sc {
		_, err := fmt.Fprint(w, v.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// Run はfinder.Selectに渡すItemsを作る
// params:
//  - ctx: context
// returns:
//  - finder.Items:
//  - error:
func (sc SourceCollection) Run(ctx context.Context) (finder.Items, error) {
	newline := []byte("\n")
	if runtime.GOOS == "window" {
		newline = []byte("\r\n")
	} else if runtime.GOOS == "darwin" {
		newline = []byte("\r")
	}

	// finder.Itemを受け取るチャンネル
	listener := make(chan finder.Item, 20)
	// errorを受け取るチャンネル
	errCh := make(chan error, 1)
	defer close(errCh)

	// 待機用
	var wg sync.WaitGroup
	wg.Add(len(sc))

	// CustomSourceを起動
	for _, v := range sc {
		go func(cs CustomSource) {
			err := cs.run(listener, newline)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(v)
	}

	// すべてのCustomSourceが終了したらlistenerを閉じる
	go func() {
		wg.Wait()
		close(listener)
	}()

	rt := finder.NewItems()
	for {
		select {
		case <-ctx.Done():
			return nil, xerrors.New("canceled by user")
		case err := <-errCh:
			return nil, err
		case item, ok := <-listener:
			if !ok {
				return rt, nil
			} else {
				rt.Add(item.Key, item.Value)
			}
		}
	}

}

// FindBeginColumn はSourceCollectionの中から name と一致するCustomSourceのBeginColumnを返す
// params:
//  - name: 探したいCustomSourceの名前
// returns:
//  - int: BeginColum
//  - error: 見つからなかったときにerror
func (sc SourceCollection) FindBeginColumn(name string) (int, error) {
	// 一度しか呼ばれないし線形探索でいいかｗ
	for _, v := range sc {
		if v.Name == name {
			return v.BeginColum, nil
		}
	}

	return -1, xerrors.Errorf("%s does not exists", name)
}
