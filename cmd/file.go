package cmd

import (
	"bufio"
	"os"
)

// Fileに追記する
func AppendRecord(path, file string) {
	// なかったら作る
	if _, err := os.Stat(file); err != nil {
	_:
		os.Create(file)
	}else{
		path += "\n"
	}

	f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err := w.WriteString(path)
	if err != nil {
		Fatal(err)
	}
	err = w.Flush()
	if err != nil {
		Fatal(err)
	}
}
