package cmd

import (
	"log"
	"os"
)

func Fatal(v ...interface{}) {
	if config.NoOutput {
		os.Exit(1)
	}

	log.Fatal(v...)
}

func Print(v ...interface{}) {
	if config.NoOutput {
		return
	}

	log.Print(v...)
}
