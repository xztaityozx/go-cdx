package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAllFile(t *testing.T) {
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

	t.Run("001_writeInputFile", func(t *testing.T) {
		f := writeInputFile(&rcd)

		if _, err := os.Stat(f); err != nil {
			t.Fatal(err)
		}

		if b, err := ioutil.ReadFile(f); err != nil {
			t.Fatal(err)
		} else {
			expect := []byte("[   1]\tABC\t/path/to/ABC\n[   2]\tDEF\t/path/to/DEF\n[   3]\tGHI\t/path/to/GHI")
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

	t.Run("003_getRecordsFromCommand", func(t *testing.T) {
		command := "/bin/ls / -1"
		out, _ := exec.Command("bash", "-c", command).Output()

		actuals := getRecordsFromCommand(command)

		for i, v := range strings.Split(string(out), "\n") {
			if !(*actuals)[i].Compare(Record{Number: i, Path: v}) {
				t.Fatal("Unexpected result getRecordsFromCommand")
			}
		}
	})
}
