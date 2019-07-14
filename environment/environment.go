package environment

import (
	"os"
	"runtime"

)

// Environment は環境依存を扱う
type Environment struct {
	// 改行文字列
	NewLine []byte
	// デフォルトシェル
	DefaultShell string
	// dev null
	DevNull string
}

// NewEnvironment はEnvironmentを作る
// returns:
//  - Environment:
func NewEnvironment() Environment {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell="sh"
	}
	if runtime.GOOS == "windows" {
		return Environment{NewLine:[]byte("\r\n"), DefaultShell:"powershell", DevNull:"$null"}
	} else if runtime.GOOS == "darwin" {
		return Environment{NewLine:[]byte("\r"), DefaultShell:shell, DevNull:"/dev/null"}
	} else {
		return Environment{NewLine:[]byte("\n"), DefaultShell:shell, DevNull:"/dev/null"}
	}
}
