package cd

import (
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/go-cdx/config"
	"testing"
)

func TestNew(t *testing.T) {
	expect := Cd{
		cfg : config.Config{
			BookmarkFile:"bookmark",
			HistoryFile:"history",
			Make:false,
			NoOutput:true,
			Source:[]config.CdxSource{},
			FuzzyFinder: config.FuzzyFinder{Command:"fzf", Options:[]string{}},
		},
		candidate: []string{"/candidate/to/candidate"},
	}
	actual := New(expect.cfg, []string{"/candidate/to/candidate"})

	assert.Equal(t, expect, actual)
}
