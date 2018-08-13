package cmd

import (
	"fmt"
	"log"
	"os"
)

func Fatal(v ...interface{}) {
	if config.NoOutput {
		fmt.Print("return 1")
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
