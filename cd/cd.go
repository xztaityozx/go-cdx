package cd

import (
	"github.com/mitchellh/go-homedir"
	"golang.org/x/xerrors"
	"os"
)

type Cd struct {
	command string
	path string
	noOutput bool
}

// New はCdをかえす
// params:
//  - cmd: /path/to/cd-command
//  - p: /path/to/destination
//  - no: 出力するかどうか
func New(cmd, p string, no bool) Cd {

	// expand path
	p, _ = homedir.Expand(p)

	return Cd{cmd,p,no}
}

// TryMakeDir は移動先のディレクトリがない場合ディレクトリを作ろうとする
// params:
//  - allow: そもそもディレクトリの作成を許可するかどうか
// returns:
//  - error: キャンセルした場合にerror
func (c Cd) TryMakeDir(allow bool) error {

	if _, err := os.Stat(c.path); err == nil {
		// 移動先がすでにあるので何もせず終了
		return nil
	} else {
		// cdx のmakeを許可しているかどうか
		if !allow {
			return xerrors.New("[cdx] make directory not allow")
		}


	}
}

