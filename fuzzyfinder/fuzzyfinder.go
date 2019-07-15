package fuzzyfinder

import (
	"context"
	"strings"

	"github.com/xztaityozx/go-cdx/environment"

	"github.com/b4b4r07/go-finder"
	"github.com/b4b4r07/go-finder/source"
	"github.com/xztaityozx/go-cdx/customsource"
)

type (
	FuzzyFinder struct {
		Path    string   `yaml:"path"`
		Options []string `yaml:"options"`
	}
)

// YesNo はyesとnoをFuzzyFinderで選択し、yesならtrue,noならfalseを返す
// returns:
//  - bool: res
//  - error:
func (ff FuzzyFinder) YesNo() (bool, error) {
	f, err := finder.New(append([]string{ff.Path}, ff.Options...)...)
	if err != nil {
		return false, err
	}

	f.Read(source.Slice([]string{"yes", "no"}))
	items, err := f.Run()

	return items[0] == "yes", err
}

// Run はFuzzyFinderを使ってパスを選択する
// params:
//  - ctx: context
//  - sc: 入力のCustomSource
// returns:
//  - string: パス
//  - error:
func (ff FuzzyFinder) Run(ctx context.Context, sc customsource.SourceCollection, env environment.Environment) (string, error) {
	f, err := finder.New(append([]string{ff.Path}, ff.Options...)...)

	if err != nil {
		return "", err
	}

	source, err := sc.Run(ctx, env)
	if err != nil {
		return "", err
	}
	res, err := f.Select(source)
	if err != nil {
		return "", err
	}

	return strings.Join(res[0].([]string), " "), nil
}
