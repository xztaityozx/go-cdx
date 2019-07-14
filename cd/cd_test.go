package cd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	home, _ := homedir.Dir()
	data := []struct{
		expect Cd
		cmd string
		p string
		no bool
	}{
		{Cd{"cd", "/path/to/A", true}, "cd", "/path/to/A", true},
		{Cd{"pushd", "/path/to/B", true}, "pushd", "/path/to/B", true},
		{Cd{"cd", "/path/to/A", false}, "cd", "/path/to/A", false},
		{Cd{"pushd", "/path/to/A", false}, "pushd", "/path/to/A", false},
		{Cd{"pushd", filepath.Join(home,"testdir"), false}, "pushd", "~/testdir", false},
	}

	for i, v := range data {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, v.expect, New(v.cmd,v.p,v.no))
		})
	}
}
