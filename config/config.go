package config

import (
	"runtime"
)

type (
	Config struct {
		Make         bool			`yaml:"make"`
		NoOutput     bool			`yaml:"noOutput"`
		Source       []CdxSource	`yaml:"source"`
		HistoryFile  string			`yaml:"history"`
		BookmarkFile string			`yaml:"bookmark"`
		FuzzyFinder  FuzzyFinder	`yaml:"fuzzyfinder"`
	}
)

var devNull = ""

func DevNull() string {
	if len(devNull) == 0 {
		if runtime.GOOS == "windows" {
			devNull = " > $null"
		} else {
			devNull = " > /dev/null"
		}
	}
	return devNull
}

var shell = ""

func DefaultShell() string {
	if len(shell) == 0 {
		if runtime.GOOS == "windows" {
			shell = "powershell.exe"
		} else {
			shell = "/bin/sh"
		}
	}
	return shell
}

var exitCommand = ""

func ExitCommand() string {
	if len(exitCommand)	== 0 {
		if runtime.GOOS == "windows" {
			exitCommand = "throw '[cdx] failed'"
		} else {
			exitCommand = "return 1"
		}
	}
	return exitCommand
}
