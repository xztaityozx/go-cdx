package fuzzyfinder

import (
	"context"
	"strings"

	"github.com/b4b4r07/go-finder"
	"github.com/xztaityozx/go-cdx/customsource"
)

type (
	FuzzyFinder struct {
		Path    string
		Options []string
	}
)

// Run はFuzzyFinderを使ってパスを選択する
// params:
//  - ctx: context
//  - sc: 入力のCustomSource
// returns:
//  - string: パス
//  - error:
func (ff FuzzyFinder) Run(ctx context.Context, sc customsource.SourceCollection) (string, error) {
	f, err := finder.New(append([]string{ff.Path}, ff.Options...)...)

	if err != nil {
		return "", err
	}

	source, err := sc.Run(ctx)
	if err != nil {
		return "", err
	}
	res, err := f.Select(source)
	if err != nil {
		return "", err
	}

	return strings.Join(res[0].([]string), " "), nil
}
