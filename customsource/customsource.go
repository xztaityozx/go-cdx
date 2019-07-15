package customsource

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xztaityozx/go-cdx/environment"

	"github.com/b4b4r07/go-finder"
)

// CustomSource
type CustomSource struct {
	// Name of CustomSource
	Name string `yaml:"name"`
	// Alias for this CustomSource
	SubName string `yaml:"alias"`
	// コマンド文字列
	Command string `yaml:"command"`
	// cdに渡したいパスが始まるカラムの先頭列
	BeginColumn int `yaml:"begin"`
}

func (cs CustomSource) Alias() rune {
	return rune(cs.SubName[0])
}

func (cs CustomSource) String() string {
	return fmt.Sprintf("%s\t%c\t%d\t%s",
		cs.Name,
		cs.Alias(),
		cs.BeginColumn,
		cs.Command,
	)
}

// run はCustomSourceに登録されたコマンドを実行して引数のoutに書き込む
// params:
//  - out: 書き込み先
//  - newline: 改行文字列
//  - env: Environment
// returns:
//  - error: コマンドの実行結果
func (cs CustomSource) run(listener chan<- finder.Item, env environment.Environment) error {

	cmd := exec.Command(env.DefaultShell, "-c", cs.Command)

	cmd.Stdin = os.Stdin
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	scan := bufio.NewScanner(stdout)
	for scan.Scan() {
		text := scan.Text()
		listener <- finder.Item{Key: fmt.Sprintf("[%s]\t%s", cs.Name, text), Value: strings.Fields(text)[cs.BeginColumn:]}
	}

	return cmd.Wait()
}
