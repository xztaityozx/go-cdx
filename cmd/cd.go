package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getCdCommand(p string) (string, error) {
	if _, err := os.Stat(p); err != nil {
		if config.Make {
			stdin := os.Stdin
			defer stdin.Close()
			stderr := os.Stderr
			defer stderr.Close()

			in := bufio.NewScanner(stdin)
			out := bufio.NewWriter(stderr)

			// --makeの質問部分
			out.Write([]byte(fmt.Sprintf("%s could not found. Make?(y/n)\n>>> ", p)))
			out.Flush()

			if in.Scan() {
				res := in.Text()
				if res == "y" || res == "yes" {
					// Directoryを作る
					if mkE := os.MkdirAll(p, 0755); mkE != nil {
						return "", mkE
					}
				} else {
					return "", errors.New("go-cdx : canceled mkdir")
				}
			}
		} else {
			return "", err
		}
	}

	// pathにスペースが入ってた時のエスケープ
	p = strings.Replace(p, " ", "\\ ", -1)

	return constructCdCommand(p), nil
}

// 実際に実行するcdコマンドを作る
func constructCdCommand(p string) string {
	base := fmt.Sprintf("%s %s", config.Command, p)
	if config.NoOutput {
		base = fmt.Sprintf("%s > /dev/null", base)
	}

	return strings.Trim(base, "\n")
}

func GetDestination(args []string) (string,string,error) {

	if !useHistory && !useBookmark && len(customSource) == 0 {
		// 履歴もブックマークもCSも使わない=>普通のパス指定
			p, _ := os.Getwd()
			if len(args) != 0 {
				p = strings.Replace(args[0], "~", os.Getenv("HOME"), 1)
				p, _ = filepath.Abs(p)
			}
		return p, "",nil
	}

	// ここからCSの領域

	var src []CustomSource
	var begins []int
	if useBookmark {
		src= append(src, BookmarkCustomSource)
		begins = append(begins, BookmarkCustomSource.BeginColumn)
	}
	if useHistory {
		src = append(src, HistoryCustomSource)
		begins = append(begins, HistoryCustomSource.BeginColumn)
	}


	m := make(map[rune]bool)

	// 複数のCustomSourceをまとめる
	for _, value := range []rune(customSource) {
		cs, err := getCustomSource(value)
		if err != nil {
			Fatal(err)
		}

		// 重複は取り除く
		if m[cs.SubName] {
			continue
		}

		m[cs.SubName]=true
		begins = append(begins, cs.BeginColumn)
		src= append(src, cs)
	}

	// CustomSourceをBuild
	c := BuildCustomSource(src...)


	path, err := c.GetPathWithFinder(begins)
	if err != nil {
		return "","",errors.New("cdx: Finder canceled")
	}

	return path, c.Action,nil
}
