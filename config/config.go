package config

import (
	"runtime"
	"github.com/xztaityozx/go-cdx/source"
)

type (
	Config struct {
		Make bool
		NoOutput bool
		Source []source.Source
		HistoryFile string
		BookmarkFile string
		FuzzyFinder
	}
)

var devNull = ""
func DevNull() string {
	if len(devNull) == 0 {
		if runtime.GOOS=="windows" {
			devNull=">$null"
		} else {
			devNull=">/dev/null"
		}
	}
	return devNull
}

var shell = ""
func DefaultShell() string {
	if len(shell) == 0 {
		if runtime.GOOS == "windows" {
			shell ="powershell.exe"
		} else {
			shell ="/bin/sh"
		}
	}
	return shell
}