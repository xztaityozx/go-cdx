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
touch %s
touch %s
`, config.BinaryPath, config.BookMarkFile, config.HistoryFile)

		actual := getInitText()

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})

}
