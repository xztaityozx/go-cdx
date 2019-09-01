package cd

import (
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/go-cdx/config"
	"github.com/xztaityozx/go-cdx/source"
	"testing"
)

func TestNew(t *testing.T) {
	expect := Cd{
		cfg : config.Config{
			BookmarkFile:"bookmark",
			HistoryFile:"history",
			Make:false,
			NoOutput:true,
			Source:[]source.Source{},
			FuzzyFinder: config.FuzzyFinder{Command:"fzf", Options:[]string{}},
		},
		dst: []string{"/dst/to/dst"},
	}
	actual := New(expect.cfg, []string{"/dst/to/dst"})

	assert.Equal(t, expect, actual)
}
