package cmd

import (
	"fmt"
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

func (act *CdxAction) Run() error {

	com := fmt.Sprintf("cd %s && %s", act.Destnation, act.Command)

	command := exec.Command("bash", "-c", com)

	if b, err := command.Output(); err != nil {
		return err
	} else {
		act.Output = b
	}
	return nil
}

func (act CdxAction) Print() {
	os.Stderr.Write(act.Output)
	os.Stderr.Close()
}
