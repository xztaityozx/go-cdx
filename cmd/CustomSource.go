package cmd

import (
	"errors"
	"fmt"
	"github.com/mattn/go-pipeline"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CustomSource struct {
	Name         string
	Command      string
	AfterCommand string
	Action       string
}

func printCustomSources() {
	os.Stderr.WriteString("cdx -c [name]\n\n[name]\t\t[command]\n")
	for _, v := range config.CustomSource {
		os.Stderr.WriteString(v.ToString())
	}
	os.Stderr.Close()
	os.Exit(0)
}

func (cs CustomSource) ToString() string {
	return fmt.Sprintf("%s\t\t%s\n", cs.Name, cs.Command)
}

func getCustomSource(name string) (string, error) {
	for _, v := range config.CustomSource {
		if v.Name == name {
			return v.Command, nil
		}
	}
	return "", errors.New(fmt.Sprint(name, "could not found in custom sorce"))
}

var HistoryCustomSource = CustomSource{
	Name:         "history",
	Command:      fmt.Sprintf("cat %s", config.HistoryFile),
	Action:       fmt.Sprintf("pwd >> %s", config.HistoryFile),
	AfterCommand: `awk '{$1=""; print}'`,
}

var BookmarkCustomSource = CustomSource{
	Name:         "bookmark",
	Command:      fmt.Sprintf("cat %s", config.BookMarkFile),
	AfterCommand: `awk '{$1=""; print}'`,
}

// TODO : ここの実装がうんこすぎる
// 複数のCustomSourceを連結する
func BuildMultipleCustomSource(list []CustomSource) (CustomSource, []int, []string) {

	command := "cat"
	action := ""
	var aftercommand []string
	var lineCount []int
	for _, v := range list {
		// Command : catへ放り込む
		command = fmt.Sprintf("%s <(%s)", command, v.Command)

		out, _ := pipeline.Output(
			[]string{"bash", "-c", v.Command},
			[]string{"wc", "-l"},
		)

		// 予期される出力
		lc, _ := strconv.Atoi(strings.Trim(string(out), "\n"))
		lineCount = append(lineCount, lc)

		// aftercommand : Sliceにいれる
		aftercommand = append(aftercommand, fmt.Sprintf("%s;%s", aftercommand, v.AfterCommand))
		// action : ;でつなぐ
		action = fmt.Sprintf("%s;%s", action, v.Action)
	}

	return CustomSource{
		Command:      command,
		AfterCommand: "",
		Action:       action,
		Name:         "multiple",
	}, lineCount, aftercommand
}

func (cs CustomSource) Equals(second CustomSource) bool {
	return cs.Action == second.Action &&
		cs.AfterCommand == second.AfterCommand &&
		cs.Command == second.Command &&
		cs.Name == second.Name
}

// FuzzyFinderをつかってパス取得。選択した行数を返す
func (cs CustomSource) GetPathWithFinder() (string, int) {
	b, err := pipeline.Output(
		[]string{"bash", "-c", cs.Command},
		config.FuzzyFinder.GetCommand(),
	)
	if err != nil {
		Fatal(err)
		os.Exit(1)
	}

	// Action
	err = exec.Command("bash", "-c", cs.Action).Run()
	if err != nil {
		Fatal(err)
	}

	splited := strings.Split(strings.Trim(string(b), "\n"), " ")
	lineNum, _ := strconv.Atoi(splited[0][1 : len(splited)-1])

	return splited[1], lineNum
}

//
