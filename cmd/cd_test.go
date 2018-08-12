package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAllCd(t *testing.T) {
	workdir := filepath.Join(os.Getenv("HOME"), "WorkSpace", "CDX")
	config = Config{
		BookMarkFile: filepath.Join(workdir, "bookmark.json"),
		HistoryFile:  filepath.Join(workdir, "history.json"),
		Command:      "echo",
		Make:         false,
		NoOutput:     false,
		UseSSH:       false,
	}

	t.Run("000_Prepare", func(t *testing.T) {
		os.MkdirAll(workdir, 0777)
	})

	t.Run("001_constructCdCommand", func(t *testing.T) {
		expect := "echo ABC"
		actual := constructCdCommand("ABC")

		if expect != actual {
			t.Fail()
		}

		expect = "echo ABC > /dev/null"
		config.NoOutput = true
		actual = constructCdCommand("ABC")
		if expect != actual {
			t.Fail()
		}

		config.NoOutput = false

	})

	t.Run("002_cd", func(t *testing.T) {
		stdin := bytes.NewBufferString("y")
		stderr := new(bytes.Buffer)
		actual, err := getCdCommand(workdir, stderr, stdin)

		if err != nil {
			t.Fatal(err)
		}

		expect := constructCdCommand(workdir)

		if expect != actual {
			t.Fatal(actual)
		}
	})
	t.Run("003_cd_Make_No", func(t *testing.T) {
		stdin := bytes.NewBufferString("n")
		stderr := new(bytes.Buffer)

		config.Make = true

		p := filepath.Join(workdir, "TEST_MAKE1")

		_, err := getCdCommand(p, stderr, stdin)

		if err == nil {
			t.Fatal("Unexpected Success")
		}

		out := []byte(fmt.Sprintf("%s could not found. Make?(y/n)\n>>> ", p))

		if bytes.Compare(out, stderr.Bytes()) != 0 {
			t.Fatal("Unexpected result stderr", string(stderr.Bytes()))
		}
	})
	t.Run("004_cd_Make_Yes", func(t *testing.T) {
		stdin := bytes.NewBufferString("y")
		stderr := new(bytes.Buffer)

		config.Make = true

		p := filepath.Join(workdir, "TEST_MAKE1")

		actual, err := getCdCommand(p, stderr, stdin)

		if err != nil {
			t.Fatal(err)
		}

		expect := constructCdCommand(p)
		out := []byte(fmt.Sprintf("%s could not found. Make?(y/n)\n>>> ", p))

		if bytes.Compare(out, stderr.Bytes()) != 0 {
			t.Fatal("Unexpected result stderr", string(stderr.Bytes()))
		}

		if expect != actual {
			t.Fatal(actual)
		}
	})

	os.RemoveAll(workdir)
}
