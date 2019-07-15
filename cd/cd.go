package cd

import (
	"context"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/xztaityozx/go-cdx/environment"
	"github.com/xztaityozx/go-cdx/fuzzyfinder"
	"golang.org/x/xerrors"
)

type Cd struct {
	command  string
	path     string
	noOutput bool
	allow    bool
	ff       fuzzyfinder.FuzzyFinder
}

// New はCdをかえす
// params:
//  - cmd: /path/to/cd-command
//  - p: /path/to/destination
//  - no: 出力するかどうか
func New(cmd, p string, no, al bool, ff fuzzyfinder.FuzzyFinder) Cd {

	// expand path
	p, _ = homedir.Expand(p)

	return Cd{cmd, p, no, al, ff}
}

// tryMakeDir は移動先のディレクトリがない場合ディレクトリを作ろうとする
// returns:
//  - error: キャンセルした場合にerror
func (c Cd) tryMakeDir() error {

	if _, err := os.Stat(c.path); err == nil {
		// 移動先がすでにあるので何もせず終了
		return nil
	} else {
		// cdx のmakeを許可しているかどうか
		if !c.allow {
			return xerrors.New("[cdx] make directory not allow")
		}

		if res, err := c.ff.YesNo(); err != nil || !res {
			return xerrors.New("[cdx] canceled")
		}

		if err := os.MkdirAll(c.path, 0755); err != nil {
			return xerrors.Errorf("[cdx] failed make directory: %w", err)
		}
	}
	return nil
}

// BuildCommand はevalするコマンドを生成する
// params:
//  - ctx: context
//  - env: environment
// returns:
//  - string: command string
func (c Cd) BuildCommand(ctx context.Context, env environment.Environment) string {

	suffix := ""
	if c.noOutput {
		suffix = fmt.Sprint(" > ", env.DevNull)
	}

	return fmt.Sprintf(`%s%s`, func() string {
		if err := c.tryMakeDir(); err != nil {
			return `exit 1`
		}

		return fmt.Sprintf(`%s '%s'`, c.command, c.path)
	}(), suffix)
}
