package config

import (
	"github.com/b4b4r07/go-finder"
	"github.com/b4b4r07/go-finder/source"
)

type FuzzyFinder struct {
	Command string
	Options []string
}

// YesNo はyesとnoをFuzzyFinderで選択し、yesならtrue,noならfalseを返す
// returns:
//  - bool: res
//  - error:
func (ff FuzzyFinder) YesNo() (bool, error) {
	f, err := finder.New(append([]string{ff.Command}, ff.Options...)...)
	if err != nil {
		return false, err
	}

	f.Read(source.Slice([]string{"yes", "no"}))
	items, err := f.Run()

	return items[0] == "yes", err
}