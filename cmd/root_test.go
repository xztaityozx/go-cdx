package cmd

import (
	"fmt"
	"testing"
)

func TestAllRoot(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		expect := fmt.Sprintf(`cdx(){
	eval "$(go-cdx $@)"
}
touch %s
touch %s
`, config.BookMarkFile, config.HistoryFile)

		actual := getInitTitle()

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})

}
