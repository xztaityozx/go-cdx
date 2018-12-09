package cmd

import (
	"fmt"
	"log"
)

func Fatal(v ...interface{}) {
	fmt.Print("return 1")

	log.Fatal(v...)
}


func Print(v ...interface{}) {
	if config.NoOutput {
		return
	}

	log.Print(v...)
}
