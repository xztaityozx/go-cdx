package cmd

import (
	"fmt"
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
	BinaryPath   string
}

type FuzzyFinder struct {
	CommandPath string
	Options     []string
}

type Record struct {
	Number int
	Path   string
}

func (rcd Record) ToString() string {
	return fmt.Sprintf("[%4d]\t%s", rcd.Number, rcd.Path)
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

type Task struct {
	History       bool
	Bookmark      bool
	CustomSources []CustomSource
}

func NewTask(h, b bool, cs []string) Task {
	rt := Task{
		History:  h,
		Bookmark: b,
	}

	return rt
}
