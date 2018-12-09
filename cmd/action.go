package cmd

import (
	"os"
	"os/exec"
)

type CdxAction struct {
	Destnation string
	Command    string
	Output     []byte
}

func NewAction(command string, dst string) CdxAction {
	return CdxAction{
		Command:    command,
		Output:     []byte(""),
		Destnation: dst,
	}
}

func (act *CdxAction) Run() {
	err := os.Chdir(act.Destnation)
	if err != nil {
		Fatal(err)
	}

	command := exec.Command("bash","-c",act.Command)
	act.Output, err = command.Output()
	if err != nil {
		Fatal(err)
	}
}

func (act CdxAction) Print() {

	_:os.Stderr.Write(act.Output)
	_:os.Stderr.Close()
}
