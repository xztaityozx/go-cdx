package source

import (
	//"context"
	//ffSource "github.com/b4b4r07/go-finder/source"
	//"github.com/b4b4r07/go-finder"
)

type  (
	Source struct {
		Name string `yaml:"name"`
		Alias rune	`yaml:"alias"`
		Command string `yaml:"command"`
		SkipColumn int `yaml:"skip"`
	}

	Collection []Source
)

//func (s Source) do(ctx context.Context) ffSource.Source {
//
//}

// toStringList is build string list for fuzzy finder
// return:
//   - []string: result
func (c Collection) toStringList() []string {
	var rt []string

	return rt
}

