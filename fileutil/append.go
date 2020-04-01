package fileutil

import (
	"fmt"
	"os"
)

func Append(path, text string) error {
	fp, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer fp.Close()

  _, err = fp.WriteString(text)
	return err
}

var newline string

func AppendLine(path, text string) error {
	fp, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer fp.Close()

  _, err = fmt.Fprintln(fp, text)
  return err
}
