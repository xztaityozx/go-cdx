package fileutil

import (
	"os"
	"runtime"
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
	if len(newline) == 0 {
		newline = func() string {
			if runtime.GOOS == "windows" {
				return "\r\n"
			} else if runtime.GOOS == "darwin" {
				return "\r"
			} else {
				return "\n"
			}
		}()
	}
	return Append(path, text+newline)
}
