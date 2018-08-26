package cmd

import (
	"bytes"
	"io"
	"io/ioutil"

	pipeline "github.com/mattn/go-pipeline"
)

// r : io.Reader (expect os.Stdin)
func readFromStdin(r io.Reader) []byte {
	rt, err := ioutil.ReadAll(r)
	if err != nil {
		Fatal(err)
	}
	return bytes.Trim(rt, "\n")
}

func sendToFuzzyFinder(b []byte) string {
	out, err := pipeline.Output(
		[]string{"echo", string(b)},
		config.FuzzyFinder.GetCommand(),
		[]string{"awk", "{print $NF}"},
	)
	if err != nil {
		Fatal(err)
	}

	return string(bytes.Trim(out, "\n"))
}

func getPathFromStdin(r io.Reader) string {
	b := readFromStdin(r)
	lines := len(bytes.Split(b, []byte("\n")))

	if lines == 1 {
		return string(b)
	}

	return sendToFuzzyFinder(b)
}
