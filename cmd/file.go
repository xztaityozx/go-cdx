package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-pipeline"
)

func writeInputFile(stream *[]Record) string {
	tmp, _ := ioutil.TempFile("", "cdxInput")
	defer tmp.Close()

	w := bufio.NewWriter(tmp)
	defer w.Flush()

	res := []string{}

	for _, v := range *stream {
		res = append(res, v.ToString())
	}

	if _, err := w.WriteString(strings.Join(res, "\n")); err != nil {
		Fatal(err)
	}

	return tmp.Name()
}

func getRecordsFromFile(f string) *[]Record {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		Fatal(err)
	}

	var res []Record
	if err := json.Unmarshal(b, &res); err != nil {
		Fatal(err)
	}

	return &res
}

func getPathWithFinderFromFile(p string) string {
	in := writeInputFile(getRecordsFromFile(p))
	b, err := pipeline.Output(
		[]string{"cat", in},
		config.FuzzyFinder.GetCommand(),
		[]string{"awk", "{print $NF}"},
	)
	if err != nil {
		Fatal(err)
		os.Exit(1)
	}

	return strings.Trim(string(b), "\n")
}

func getPathWithFinderFromCommand(command string) string {
	b, err := pipeline.Output(
		[]string{"bash", "-c", command},
		config.FuzzyFinder.GetCommand(),
		[]string{"awk", "{print $NF}"},
	)

	if err != nil {
		Fatal(err)
		os.Exit(1)
	}

	return strings.Trim(string(b), "\n")
}

func getPathWithFinder() (string, error) {
	if useBookmark {
		return getPathWithFinderFromFile(config.BookMarkFile), nil
	} else if useHistory {
		return getPathWithFinderFromFile(config.HistoryFile), nil
	} else if len(customSource) != 0 {
		name, err := getCustomSource(customSource)
		if err != nil {
			return "", err
		}
		return getPathWithFinderFromCommand(name), nil
	}
	return "", errors.New("")
}

func AppendRecord(p string, target string) {
	res := *getRecordsFromFile(target)
	n := len(res) + 1
	res = append(res, Record{
		Number: n,
		Path:   p,
	})

	b, err := json.Marshal(res)
	if err != nil {
		Fatal(err)
	}
	if err := ioutil.WriteFile(target, b, 0644); err != nil {
		Fatal(err)
	}
}

func GetDestination(args []string) string {
	if fuz, err := getPathWithFinder(); err == nil {
		return fuz
	} else {
		p, _ := os.Getwd()
		if len(args) != 0 {
			p = strings.Replace(args[0], "~", os.Getenv("HOME"), 1)
			p, _ = filepath.Abs(p)
		}
		return p
	}
}

func TryCreateFiles(p string) error {
	if _, err := os.Stat(p); err != nil {
		return ioutil.WriteFile(p, []byte("[]"), 0644)
	}
	return nil
}
