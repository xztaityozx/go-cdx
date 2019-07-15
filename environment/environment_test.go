package environment

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_New(t *testing.T) {
	var expect Environment
	if runtime.GOOS == "windows" {
		expect = Environment{[]byte("\r\n"), "powershell", "$null"}
	} else {
		shell := os.Getenv("SHELL")
		if len(shell) == 0 {
			shell = "sh"
		}

		if runtime.GOOS == "darwin" {
			expect = Environment{[]byte("\r"), shell, "/dev/null"}
		} else {
			expect = Environment{[]byte("\n"), shell, "/dev/null"}
		}
	}

	assert.Equal(t, expect, NewEnvironment())

}
