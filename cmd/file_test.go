package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAllFile(t *testing.T) {
	workdir := filepath.Join(os.Getenv("HOME"), "WorkSpace", "CDX")
	config = Config{
		BookMarkFile: filepath.Join(workdir, "bookmark.json"),
		HistoryFile:  filepath.Join(workdir, "history.json"),
		FuzzyFinder: FuzzyFinder{
			CommandPath: "/usr/bin/head",
			Options:     []string{"-n1"},
		},
		CustomSource: []CustomSource{
			CustomSource{
				Command: "/custom/command",
				SubName: 'c',
			},
		},
		Command:  "echo",
		Make:     false,
		NoOutput: false,
	}

	t.Run("000_Prepare", func(t *testing.T) {
		if err := os.MkdirAll(workdir, 0777); err != nil {
			t.Fatal(err)
		}
	})


	t.Run("005_getCustomSource", func(t *testing.T) {
		expect := config.CustomSource[0].Command
		actual, err := getCustomSource(config.CustomSource[0].SubName)
		if err != nil {
			t.Fatal(err)
		}
		if expect != actual.Command {
			t.Fatal(actual, "is not", expect)
		}

		if _, err := getCustomSource('A'); err == nil {
			t.Fatal("Unexpected result getCustomSorce")
		}

	})

	t.Run("006_AppendRecord", func(t *testing.T) {
		expect := "/path/to"
		path := filepath.Join(os.Getenv("HOME"),"workdir")
		err := os.MkdirAll(path,0777)
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(path,"appendfile")
		AppendRecord(expect,target)

		actual, _ := ioutil.ReadFile(target)
		if expect != string(actual) {
			t.Fatal(expect,"is not",string(actual))
		}
		os.Remove(target)

	})

}
