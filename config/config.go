package config

import (
	"bufio"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/xztaityozx/go-cdx/customsource"
	"github.com/xztaityozx/go-cdx/fuzzyfinder"
)

type (
	Config struct {
		File          File                        `yaml:"file"`
		Command       string                      `yaml:"command"`
		NoOutput      bool                        `yaml:"noOutput"`
		Make          bool                        `yaml:"make"`
		FuzzyFinder   fuzzyfinder.FuzzyFinder     `yaml:"fuzzyfinder"`
		CustomSources []customsource.CustomSource `yaml:"source"`
	}

	File struct {
		History  string `yaml:"history"`
		BookMark string `yaml:"bookmark"`
	}
)

func appendTo(t, p string) error {
	fp, err := os.OpenFile(t, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	w := bufio.NewWriter(fp)

	_, err = w.WriteString(p)
	return err
}

func (f File) AppendBookmark(p string) error {
	b, _ := homedir.Expand(f.BookMark)
	return appendTo(b, p)
}

func (f File) AppendHistory(p string) error {
	h, _ := homedir.Expand(f.History)
	return appendTo(h, p)
}
