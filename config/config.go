package config

import (
	"github.com/xztaityozx/go-cdx/customsource"
	"github.com/xztaityozx/go-cdx/fuzzyfinder"
)

type (
	Config struct {
		File          File
		Command       string
		NoOutput      bool
		Make          bool
		FuzzyFinder   fuzzyfinder.FuzzyFinder
		CustomSources []customsource.CustomSource
	}

	File struct {
		History  string
		BookMark string
	}
)
