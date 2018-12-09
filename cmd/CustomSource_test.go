package cmd

import "testing"

func TestAllCustomSource(t *testing.T) {
	t.Run("001_BuildMultipleCustomSource", func(t *testing.T) {
		src := []CustomSource{
			{
				SubName: 'f',
				Command: "echo command first",
				Action:  "echo action first",
			},
			{
				SubName: 's',
				Command: "echo command second",
				Action:  "echo action second",
			},
			{
				SubName: 't',
				Command: "echo command third",
				Action:  "echo action third",
			},
		}

		actual := BuildCustomSource(src...)
		expect := CustomSource{
			SubName:     'm',
			BeginColumn: -1,
			Command:     `cat <(echo command first|awk '{printf "[1-%03d] %s\n",NR,$0}'|column -t) <(echo command second|awk '{printf "[2-%03d] %s\n",NR,$0}'|column -t) <(echo command third|awk '{printf "[3-%03d] %s\n",NR,$0}'|column -t)`,
			Action:      " echo action first; echo action second; echo action third;",
		}

		if !actual.Equals(expect) {
			t.Fatal(actual, "is not \n", expect)
		}

	})
}
