package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func getCdCommand(p string, stderr io.Writer, stdin io.Reader) (string, error) {
	if _, err := os.Stat(p); err != nil {
		if config.Make {
			in := bufio.NewScanner(stdin)
			out := bufio.NewWriter(stderr)

			out.Write([]byte(fmt.Sprintf("%s could not found. Make?(y/n)\n>>> ", p)))
			out.Flush()

			if in.Scan() {
				res := in.Text()
				if res == "y" || res == "yes" {
					// Directoryを作る
					if mkE := os.MkdirAll(p, 0755); mkE != nil {
						return "", mkE
					}
				} else {
					return "", errors.New("go-cdx : canceled mkdir")
				}
			}
		} else {
			return "", err
		}
	}

	return constructCdCommand(p), nil
}

func constructCdCommand(p string) string {
	base := fmt.Sprintf("%s %s", config.Command, p)
	if config.NoOutput {
		base = fmt.Sprintf("%s > /dev/null", base)
	}

	return strings.Trim(base, "\n")
}
