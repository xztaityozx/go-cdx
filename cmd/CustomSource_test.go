package cmd

import "testing"

func TestAllCustomSource(t *testing.T) {
	t.Run("001_BuildMultipleCustomSource", func(t *testing.T) {
		src := []CustomSource{
			{
				Name:         "first",
				Command:      "echo command first",
				AfterCommand: "echo after first",
				Action:       "echo action first",
			},
			{
				Name:         "second",
				Command:      "echo command second",
				AfterCommand: "echo after second",
				Action:       "echo action second",
			},
			{
				Name:         "third",
				Command:      "echo command third",
				AfterCommand: "echo after third",
				Action:       "echo action third",
			},
		}

		actual, lines, afters := BuildMultipleCustomSource(src)
		expect := CustomSource{
			Name:         "multiple",
			Command:      "cat <(echo command first) <(echo command second) <(echo command third)",
			AfterCommand: "",
			Action:       "action first;action second;action third",
		}

		if !actual.Equals(expect) {
			t.Fatal(actual, "is not ", expect)
		}

		if len(afters) != 3 {
			t.Fatal("error: ", afters)
		}

		if len(lines) != 3 {
			t.Fatal("error: ", lines)
		}

	})
}
