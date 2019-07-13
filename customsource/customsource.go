package customsource

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/b4b4r07/go-finder"
)

// CustomSource
type CustomSource struct {
	// Name of CustomSource
	Name string
	// Alias for this CustomSource
	Alias rune
	// コマンド文字列
	Command string
	// cdに渡したいパスが始まるカラムの先頭列
	BeginColum int
}

func (cs CustomSource) String() string {
	return fmt.Sprintf("%s\t%c\t%d\t%s",
		cs.Name,
		cs.Alias,
		cs.BeginColum,
		cs.Command,
	)
}

// run はCustomSourceに登録されたコマンドを実行して引数のoutに書き込む
// params:
//  - out: 書き込み先
//  - newline: 改行文字列
// returns:
//  - error: コマンドの実行結果
func (cs CustomSource) run(listener chan<- finder.Item, newline []byte) error {

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cs.Command)
	} else {
		cmd = exec.Command("sh", "-c", cs.Command)
	}

	cmd.Stdin = os.Stdin
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	//err = cmd.Start()
	//if err != nil {
	//return err
	//}

	ch := make(chan struct{})
	defer close(ch)
	go func() {
		scan := bufio.NewScanner(stdout)
		for scan.Scan() {
			text := scan.Text()
			listener <- finder.Item{Key: fmt.Sprintf("[%s]\t%s", cs.Name, text), Value: strings.Fields(text)[cs.BeginColum:]}
		}
		ch <- struct{}{}
	}()

	err = cmd.Run()

	// block
	<-ch
	return err
}
