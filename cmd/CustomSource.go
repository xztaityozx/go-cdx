package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-pipeline"
)

type CustomSource struct {
	Name 		string
	SubName     rune
	Command     string
	BeginColumn int
	Action      string
}

func printCustomSources() {
	_, _ = os.Stderr.WriteString("cdx -c [name]\n\n[name]\t\t[command]\n")
	for _, v := range config.CustomSource {
		_,_ = os.Stderr.WriteString(v.ToString())
	}
	_ = os.Stderr.Close()
	os.Exit(0)
}

func (cs CustomSource) ToString() string {
	return fmt.Sprintf("%s\t\t%s\n", cs.Name, cs.Command)
}

func getCustomSource(name rune) (CustomSource, error) {
	for _, v := range config.CustomSource {
		if (len(v.Name) != 0 && []rune(v.Name)[0] == name)  || v.SubName == name{
			return v, nil
		}
	}
	return CustomSource{}, errors.New(fmt.Sprint(name, "could not found in custom source"))
}

var HistoryCustomSource CustomSource
var BookmarkCustomSource CustomSource



// 複数のCustomSourceを連結する
func BuildCustomSource(list... CustomSource) CustomSource {

	command := "cat"
	action := ""

	for k, v := range list {
		command = fmt.Sprintf(`%s <(%s|awk '{printf "[%d-%%03d] %%s\n",NR,$0}'|column -t)`, command, v.Command, k+1)
		if len(v.Action) != 0 {
			action = fmt.Sprintf("%s %s;", action,v.Action)
		}
	}

	return CustomSource{
		SubName:     'm',
		Command:     command,
		BeginColumn: -1,
		Action:      action,
	}
}

func (cs CustomSource) Equals(second CustomSource) bool {
	return cs.Action == second.Action &&
		cs.BeginColumn == second.BeginColumn &&
		cs.Command == second.Command &&
		cs.SubName == second.SubName
}

// FuzzyFinderをつかってパス取得。AfterCommandのインデックスを返す
func (cs CustomSource) GetPathWithFinder(begins []int) (string,error) {

	b, err := pipeline.Output(
		[]string{"bash", "-c", cs.Command},
		config.FuzzyFinder.GetCommand(),
	)
	if err != nil {
		return "", err
	}


	// 1つ以上の空白で区切る
	field := strings.Fields(string(b))

	// BeginColumnのIndexを得る
	var idx string
	for _, value := range field[0] {
		if value=='[' {
			continue
		}
		if value == '-' {
			break
		}
		idx = fmt.Sprintf("%s%c",idx,value)
	}

	beginsIdx, _ := strconv.Atoi(idx)


	// 1 indexed => 0 indexed
	return strings.Join(field[begins[beginsIdx-1]:]," "), nil
}

