package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/xerrors"
)

func Initialize() error {
	gopath := os.Getenv("GOPATH")
	repoPath := filepath.Join(gopath, "src", "github.com", "xztaityozx", "go-cdx")

	scriptPath := filepath.Join(repoPath, "script", func() string {
		if runtime.GOOS == "windows" {
			return "cdx.ps1"
		} else {
			return "cdx.sh"
		}
	}())

	fp, err := os.Open(scriptPath)
	if err != nil {
		return xerrors.Errorf("cdx cannot found init script :%v", scriptPath)
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	for s.Scan() {
		fmt.Println(s.Text())
	}

	return nil
}
