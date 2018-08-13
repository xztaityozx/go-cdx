package cmd

import (
	"fmt"
	"testing"
)

func TestAllRoot(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		expect := fmt.Sprintf(`cdx(){
	eval "$(%s $@)"
}
[ -f %s ] || echo -n "[]" > %s
[ -f %s ] || echo -n "[]" > %s
`, config.BinaryPath, config.BookMarkFile, config.BookMarkFile, config.HistoryFile, config.HistoryFile)

		actual := getInitText()

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})

}
