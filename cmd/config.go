package cmd

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	HistoryFile  string
	BookMarkFile string
	Command      string
	NoOutput     bool
	UseSSH       bool
	Make         bool
	CustomSource []CustomSource
	FuzzyFinder  FuzzyFinder
}

type FuzzyFinder struct {
	CommandPath string
	Options     []string
}

type CustomSource struct {
	Name    string
	Command string
}

type Record struct {
	Number int
	Path   string
}

func (rcd Record) ToString() string {
	base := filepath.Base(rcd.Path)
	return fmt.Sprintf("[%4d]\t%s\t%s", rcd.Number, base, rcd.Path)
}

func (s Record) Compare(t Record) bool {
	return s.Number == t.Number && s.Path == t.Path
}

func (ff FuzzyFinder) GetCommand() []string {
	var rt []string = []string{
		ff.CommandPath,
	}

	rt = append(rt, ff.Options...)

	return rt
}
