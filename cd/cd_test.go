package cd

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/go-cdx/environment"
	"github.com/xztaityozx/go-cdx/fuzzyfinder"
)

func TestNew(t *testing.T) {
	home, _ := homedir.Dir()
	data := []struct {
		expect Cd
		cmd    string
		p      string
		no     bool
		al     bool
		ff     fuzzyfinder.FuzzyFinder
	}{
		{Cd{"cd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "cd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"cd", filepath.Join(home, "test"), true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "cd", "~/test", true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"pushd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "pushd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"pushd", filepath.Join(home, "test"), true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "pushd", "~/test", true, true, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"cd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "cd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"cd", filepath.Join(home, "test"), false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "cd", "~/test", false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"pushd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "pushd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"pushd", filepath.Join(home, "test"), false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}}, "pushd", "~/test", false, false, fuzzyfinder.FuzzyFinder{Path: "fzf", Options: []string{}}},
		{Cd{"cd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "cd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"cd", filepath.Join(home, "test"), true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "cd", "~/test", true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"pushd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "pushd", "/path/to/A", true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"pushd", filepath.Join(home, "test"), true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "pushd", "~/test", true, true, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"cd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "cd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"cd", filepath.Join(home, "test"), false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "cd", "~/test", false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"pushd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "pushd", "/path/to/A", false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
		{Cd{"pushd", filepath.Join(home, "test"), false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}}, "pushd", "~/test", false, false, fuzzyfinder.FuzzyFinder{Path: "peco", Options: []string{}}},
	}

	for i, v := range data {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, v.expect, New(v.cmd, v.p, v.no, v.al, v.ff))
		})
	}
}

func TestCd_BuildCommand(t *testing.T) {

	sha := fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().Format(time.ANSIC))))
	home := os.TempDir()
	data := []struct {
		expect string
		env    environment.Environment
		source Cd
	}{
		{expect: fmt.Sprintf(`cd '%s'`, home), env: environment.Environment{}, source: Cd{"cd", home, false, false, fuzzyfinder.FuzzyFinder{}}},
		{expect: fmt.Sprintf(`cd '%s' > /dev/null`, home), env: environment.Environment{DevNull: "/dev/null"}, source: Cd{"cd", home, true, false, fuzzyfinder.FuzzyFinder{}}},
		{expect: fmt.Sprintf(`cd '%s' > $null`, home), env: environment.Environment{DevNull: "$null"}, source: Cd{"cd", home, true, false, fuzzyfinder.FuzzyFinder{}}},
		{expect: fmt.Sprintf(`pushd '%s'`, home), env: environment.Environment{}, source: Cd{"pushd", home, false, false, fuzzyfinder.FuzzyFinder{}}},
		{expect: fmt.Sprintf(`pushd '%s' > /dev/null`, home), env: environment.Environment{DevNull: "/dev/null"}, source: Cd{"pushd", home, true, false, fuzzyfinder.FuzzyFinder{}}},
		{expect: fmt.Sprintf(`pushd '%s' > $null`, home), env: environment.Environment{DevNull: "$null"}, source: Cd{"pushd", home, true, false, fuzzyfinder.FuzzyFinder{}}},
		{expect: fmt.Sprintf(`cd '%s'`, home), env: environment.Environment{}, source: Cd{"cd", home, false, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`cd '%s' > /dev/null`, home), env: environment.Environment{DevNull: "/dev/null"}, source: Cd{"cd", home, true, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`cd '%s'`, home), env: environment.Environment{}, source: Cd{"cd", home, false, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`cd '%s' > $null`, home), env: environment.Environment{DevNull: "$null"}, source: Cd{"cd", home, true, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`pushd '%s'`, home), env: environment.Environment{}, source: Cd{"pushd", home, false, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`pushd '%s' > /dev/null`, home), env: environment.Environment{DevNull: "/dev/null"}, source: Cd{"pushd", home, true, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`pushd '%s'`, home), env: environment.Environment{}, source: Cd{"pushd", home, false, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},
		{expect: fmt.Sprintf(`pushd '%s' > $null`, home), env: environment.Environment{DevNull: "$null"}, source: Cd{"pushd", home, true, false, fuzzyfinder.FuzzyFinder{Path: "head", Options: []string{"-n1"}}}},

		{
			expect: `exit 1`,
			env:    environment.Environment{},
			source: Cd{
				"cd",
				filepath.Join(home, "test", sha),
				false,
				true,
				fuzzyfinder.FuzzyFinder{Path: "echo", Options: []string{"no"}},
			},
		},
		{
			expect: fmt.Sprintf(`cd '%s'`, filepath.Join(home, "test", sha)),
			env:    environment.Environment{},
			source: Cd{
				"cd",
				filepath.Join(home, "test", sha),
				false,
				true,
				fuzzyfinder.FuzzyFinder{Path: "echo", Options: []string{"yes"}},
			},
		},
		{
			expect: `exit 1`,
			env:    environment.Environment{},
			source: Cd{
				"cd",
				filepath.Join(home, "test", sha),
				false,
				true,
				fuzzyfinder.FuzzyFinder{Path: "echo", Options: []string{"no"}},
			},
		},
		{
			expect: fmt.Sprintf(`cd '%s'`, filepath.Join(home, "test", sha)),
			env:    environment.Environment{},
			source: Cd{
				"cd",
				filepath.Join(home, "test", sha),
				false,
				true,
				fuzzyfinder.FuzzyFinder{Path: "echo", Options: []string{"yes"}},
			},
		},
		{
			expect: `exit 1 > /dev/null`,
			env:    environment.Environment{DevNull: "/dev/null"},
			source: Cd{
				"cd",
				filepath.Join(home, "test", sha),
				true,
				true,
				fuzzyfinder.FuzzyFinder{Path: "echo", Options: []string{"no"}},
			},
		},
		{
			expect: `exit 1 > $null`,
			env:    environment.Environment{DevNull: "$null"},
			source: Cd{
				"cd",
				filepath.Join(home, "test", sha),
				true,
				true,
				fuzzyfinder.FuzzyFinder{Path: "echo", Options: []string{"no"}},
			},
		},
	}

	for i, v := range data {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, v.expect, v.source.BuildCommand(context.Background(), v.env))
			os.RemoveAll(filepath.Join(home, "test"))
		})
	}

}
