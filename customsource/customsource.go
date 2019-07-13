package customsource

import (
	"fmt"
	"github.com/mattn/go-pipeline"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
)

// CustomSource
type CustomSource struct {
	// Name of CustomSource
	Name string
	// Alias for this CustomSource
	Alias       rune
	Commands    [][]string
}

func (cs CustomSource) String() string {
	var box []string
	for _, v := range cs.Commands {
		box = append(box, strings.Join(v, " "))
	}
	return fmt.Sprintf("%s\t%c\t%s",
		cs.Name,
		cs.Alias,
		strings.Join(box, "|"),
	)
}

func PrintCustomSource(cs []CustomSource) error {
	w := tabwriter.NewWriter(os.Stderr, 0, 8, 0, '\t', 0)
	_, err := fmt.Fprint(w, "Name\tAlias\tCommands")
	if err != nil {
		return err
	}
	for _, v := range cs {
		_, err := fmt.Fprint(w, v.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// Start start commands
func (cs CustomSource) Start() ([]string, error) {
	out, err := pipeline.Output(cs.Commands...)

	sep := "\n"
	if runtime.GOOS == "windows" {
		sep = "\r\n"
	} else if runtime.GOOS == "darwin" {
		sep = "\r"
	}


	return strings.Split(string(out), sep), err
}
