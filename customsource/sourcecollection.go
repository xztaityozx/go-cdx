package customsource

import (
	"context"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/xztaityozx/go-cdx/environment"

	"github.com/b4b4r07/go-finder"
	"golang.org/x/xerrors"
)

type SourceCollection []CustomSource

// New はconfigで設定しているCustomSourceから、listの各文字がAliasと一致するもののリストを返す
// params:
//  - list: 一致させたいAliasのリスト
//  - box: 設定されてるCustomSource
// returns:
//  - SourceCollection:
func New(list string, box []CustomSource) (SourceCollection, error) {
	m := map[rune]CustomSource{}
	for _, v := range box {
		m[v.Alias] = v
	}

	var rt SourceCollection
	for _, v := range []rune(list) {
		t, ok := m[v]
		if ok {
			rt = append(rt, t)
		} else {
			return nil, xerrors.Errorf("%s not found", v)
		}
	}
	return rt, nil
}

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
//  - env: Environment
// returns:
//  - finder.Items:
//  - error:
func (sc SourceCollection) Run(ctx context.Context, env environment.Environment) (finder.Items, error) {

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
			err := cs.run(listener, env)
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
