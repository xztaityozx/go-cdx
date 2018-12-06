package cmd

import (
	"bytes"
	"encoding/json"
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
				Name:    "custom",
			},
		},
		Command:  "echo",
		Make:     false,
		NoOutput: false,
		UseSSH:   false,
	}
	rcd := []Record{
		Record{
			Number: 1,
			Path:   "/path/to/ABC",
		},
		Record{
			Number: 2,
			Path:   "/path/to/DEF",
		},
		Record{
			Number: 3,
			Path:   "/path/to/GHI",
		},
	}

	t.Run("000_Prepare", func(t *testing.T) {
		if err := os.MkdirAll(workdir, 0777); err != nil {
			t.Fatal(err)
		}
		TryCreateFiles(config.BookMarkFile)
		TryCreateFiles(config.HistoryFile)

	})

	t.Run("001_writeInputFile", func(t *testing.T) {
		f := writeInputFile(&rcd)

		if _, err := os.Stat(f); err != nil {
			t.Fatal(err)
		}

		if b, err := ioutil.ReadFile(f); err != nil {
			t.Fatal(err)
		} else {
			expect := []byte("[   1]\t/path/to/ABC\n[   2]\t/path/to/DEF\n[   3]\t/path/to/GHI")
			if bytes.Compare(expect, b) != 0 {
				t.Fatal("Unexpect file Item")
			}
		}
	})

	t.Run("002_getRecordsFromFile", func(t *testing.T) {
		b, _ := json.Marshal(rcd)
		tmp, _ := ioutil.TempFile("", "test002_FromFile")
		defer os.Remove(tmp.Name())
		tmp.Close()

		ioutil.WriteFile(tmp.Name(), b, 0644)

		for i, v := range *getRecordsFromFile(tmp.Name()) {
			if !rcd[i].Compare(v) {
				t.Fatal("Unexpected result getRecordsFromFile")
			}
		}
	})

	t.Run("003_getCdCommandWithFinderFromCommand", func(t *testing.T) {
		actual := getPathWithFinderFromCommand("echo ABC")
		expect := "ABC"

		if expect != actual {
			t.Fatal(actual, "is not", expect)
		}
	})

	t.Run("004_getCdCommandWithFinderFromFile", func(t *testing.T) {
		b, _ := json.Marshal(rcd)
		tmp, _ := ioutil.TempFile("", "test004_FromFile")
		defer os.Remove(tmp.Name())
		tmp.Close()

		ioutil.WriteFile(tmp.Name(), b, 0644)

		expect := rcd[0].Path
		actual := getPathWithFinderFromFile(tmp.Name())

		if expect != actual {
			t.Fatal(actual, "is not", expect)
		}

	})

	t.Run("005_getCustomSource", func(t *testing.T) {
		expect := config.CustomSource[0].Command
		actual, err := getCustomSource(config.CustomSource[0].Name)
		if err != nil {
			t.Fatal(err)
		}
		if expect != actual {
			t.Fatal(actual, "is not", expect)
		}

		if _, err := getCustomSource("ABC"); err == nil {
			t.Fatal("Unexpected result getCustomSorce")
		}

	})

	t.Run("006_AppendRecord", func(t *testing.T) {
		os.MkdirAll(workdir, 0777)
		ioutil.WriteFile(config.BookMarkFile, []byte("[]"), 0644)

		AppendRecord("/path/to", config.BookMarkFile)

		expect := `[{"Number":1,"Path":"/path/to"}]`
		actual, _ := ioutil.ReadFile(config.BookMarkFile)

		if expect != string(actual) {
			t.Fatal(string(actual), "is not", expect)
		}

		os.Remove(config.BookMarkFile)
	})

}
