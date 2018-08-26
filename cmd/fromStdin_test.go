package cmd

import "bytes"
import "strings"
import "testing"

func TestAllFromStdin(t *testing.T) {
	t.Run("001_readFromStdin", func(t *testing.T) {
		r := strings.NewReader("ABC\nDEF\nHIJ")
		s := readFromStdin(r)

		expect := []byte("ABC\nDEF\nHIJ")

		if bytes.Compare(expect, s) != 0 {
			t.Fatal(s, "is not", expect)
		}

	})

	t.Run("002_sendToFuzzyFinder", func(t *testing.T) {
		config.FuzzyFinder = FuzzyFinder{
			CommandPath: "head",
			Options:     []string{"-n1"},
		}
		r := strings.NewReader("ABC\nDEF\nHIJ")
		b := readFromStdin(r)

		actual := sendToFuzzyFinder(b)
		expect := "ABC"

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}

	})

	t.Run("003_getPathFromStdin_one", func(t *testing.T) {
		r := strings.NewReader("ABC")

		actual := getPathFromStdin(r)
		expect := "ABC"

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})

	t.Run("004_getPathFromStdin_multi", func(t *testing.T) {
		r := strings.NewReader("ABC\nDEF\nHIJ")

		actual := getPathFromStdin(r)
		expect := "ABC"

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}

	})

}
