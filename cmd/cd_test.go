package cmd

import (
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
		actual, err := getCdCommand(workdir)

		if err != nil {
			t.Fatal(err)
		}

		expect := constructCdCommand(workdir)

		if expect != actual {
			t.Fatal(actual)
		}
	})

	t.Run("003_GetDestination_cs_1", func(t *testing.T){
		config.CustomSource = []CustomSource{
			{
				Name:"u",
				SubName:     'u',
				BeginColumn: 2,
				Command:     "yes first second|head -n1",
				Action:      "",
			},
		}
		config.FuzzyFinder = FuzzyFinder{
			CommandPath:"cat",
		}


		customSource="u"
		actual, _, _ := GetDestination([]string{})
		expect := "second"

		if actual!=expect{
			Fatal(actual,"is not",expect)
		}
	})
	t.Run("004_GetDestination_cs_uk", func(t *testing.T){
		config.CustomSource = []CustomSource{
			{
				SubName:     'u',
				BeginColumn: 2,
				Command:     "echo command u",
				Action:      "",
			},
			{
				SubName:     'k',
				BeginColumn: 1,
				Command:     "echo command k",
				Action:      "",
			},
		}
		config.FuzzyFinder = FuzzyFinder{
			CommandPath:"head",
			Options: []string{"-n1"},
		}


		customSource="uk"
		actual, _, _ := GetDestination([]string{})
		expect := "u"

		if actual!=expect{
			Fatal(actual,"is not",expect)
		}
	})
	t.Run("005_GetDestination_cs_ku", func(t *testing.T){
		config.CustomSource = []CustomSource{
			{
				SubName:     'u',
				BeginColumn: 2,
				Command:     "echo command u",
				Action:      "",
			},
			{
				SubName:     'k',
				BeginColumn: 1,
				Command:     "echo command k",
				Action:      "",
			},
		}
		config.FuzzyFinder = FuzzyFinder{
			CommandPath:"head",
			Options: []string{"-n1"},
		}


		customSource="ku"
		actual, _,_ := GetDestination([]string{})
		expect := "command k"

		if actual!=expect{
			Fatal(actual,"is not",expect)
		}
	})
	t.Run("006_GetDestination_no_finder", func(t *testing.T){
		useBookmark=false
		useHistory=false
		customSource=""
		actual, _,_ := GetDestination([]string{"~/"})
		expect := filepath.Join(os.Getenv("HOME"))

		if actual!=expect{
			Fatal(actual,"is not",expect)
		}
	})

	t.Run("006_GetDestination_wd", func(t *testing.T){
		useBookmark=false
		useHistory=false
		customSource=""
		actual, _,_ := GetDestination([]string{})
		wd,_:=os.Getwd()
		expect := filepath.Join(wd)

		if actual!=expect{
			Fatal(actual,"is not",expect)
		}
	})

	os.RemoveAll(workdir)
}
