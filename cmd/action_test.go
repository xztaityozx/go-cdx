package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestAllAction(t *testing.T) {
	home := os.Getenv("HOME")

	t.Run("001_NewAction", func(t *testing.T) {
		a := NewAction("ABC", "DEF")
		if a.Command != "ABC" {
			t.Fatal("Unexpected result a.Command")
		}
		if a.Destnation != "DEF" {
			t.Fatal("Unexpected result a.Destnation")
		}

	})

	t.Run("002_Run", func(t *testing.T) {
		act := NewAction("echo abc", home)
		act.Run()
		if bytes.Compare(act.Output, []byte("abc\n")) != 0 {
			t.Fatal("Unexpected result : ", act.Output)
		}
	})
}
