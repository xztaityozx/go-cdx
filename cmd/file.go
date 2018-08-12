package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func getCdCommandWithFinderFromFile(p string) string {
	in := writeInputFile(getRecordsFromFile(p))
	b, err := pipeline.Output(
		[]string{"cat", in},
		config.FuzzyFinder.GetCommand(),
		[]string{"awk", "{print $NF}"},
	)
	if err != nil {
		Fatal(err)
	}

	return constructCdCommand(string(b))
}

func getCdCommandWithFinderFromCommand(command string) string {
	b, err := pipeline.Output(
		[]string{"bash", "-c", command},
		config.FuzzyFinder.GetCommand(),
		[]string{"awk", "{print $NF}"},
	)

	if err != nil {
		Fatal(err)
	}

	return constructCdCommand(string(b))
}

func getCustomSorce(name string) (string, error) {
	for _, v := range config.CustomSource {
		if v.Name == name {
			return v.Command, nil
		}
	}
	return "", errors.New(fmt.Sprint(name, "could not found in custom sorce"))
}

func getCdCommandWithFinder() (string, error) {

	if useBookmark {
		return getCdCommandWithFinderFromFile(config.BookMarkFile), nil
	} else if useHistory {
		return getCdCommandWithFinderFromFile(config.HistoryFile), nil
	} else if len(customSource) != 0 {
		name, err := getCustomSorce(customSource)
		if err != nil {
			return "", err
		}
		return getCdCommandWithFinderFromCommand(name), nil
	}
	return "", errors.New("")
}
