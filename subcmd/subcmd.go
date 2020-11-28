package subcmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

// for -p,--popd option
func Popd() (string, error) {
	return "popd", nil
}

// for --init option
func Initialize() (string, error) {
	if runtime.GOOS == "windows" {
		return fmt.Sprintln(`
function cdx() {
    begin{
        [System.Collections.Generic.List[string]]$paths=@{}
    }
    process{
        $paths.Add($_)
    }
    end{
        $command = "$($paths | go-cdx $args)";
        Invoke-Expression "$command"
    }
}`), nil
	} else {
		return fmt.Sprint(`function() cdx() {command="$(go-cdx $@)"; eval "${command}"}`), nil
	}
}

// for --completion option
func GenCompletion(cmd *cobra.Command, args []string) {
	shell, _ := cmd.Flags().GetString("completion")
	if len(shell) == 0 {
		return
	}

	if shell == "bash" {
		_ = cmd.GenBashCompletion(os.Stdout)
	} else if shell == "zsh" {
		_ = cmd.GenZshCompletion(os.Stdout)
	} else if shell == "fish" {
		_ = cmd.GenFishCompletion(os.Stdout, false)
	} else if shell == "PowerShell" {
		_ = cmd.GenPowerShellCompletion(os.Stdout)
	} else {
		logrus.Fatal(shell, " is unsupported shell name")
	}

	os.Exit(0)
}

// for --add option
func Add(bookmarkFile string) (string, error) {
	cwd, _ := os.Getwd()
	f, err := os.OpenFile(bookmarkFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return "false", err
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, cwd)
	if err != nil {
		return "false", err
	}
	return "true", nil
}

// for -R, --git-root option
func GitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-cdup")

	o, err := cmd.CombinedOutput()
	o = bytes.Trim(o, "\n")
	if err != nil {
		return "false", xerrors.New(string(o))
	}

	if len(o) == 0 {
		return fmt.Sprint("true"), nil
	} else {
		return fmt.Sprintf("pushd %s", string(o)), nil
	}
}
