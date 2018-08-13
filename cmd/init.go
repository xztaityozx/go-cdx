package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func PrintInitText() {
	if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "go-cdx"), 0777); err != nil {
		Fatal(err)
	}

	fmt.Print(getInitText())
	os.Exit(0)
}

func getInitText() string {
	return fmt.Sprintf(`cdx(){
	eval "$(%s $@)"
}
touch %s
touch %s
`, config.BinaryPath, config.BookMarkFile, config.HistoryFile)
}
