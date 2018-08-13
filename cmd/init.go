package cmd

import (
	"fmt"
)

func PrintInitText() {

}

func getInitTitle() string {
	return fmt.Sprintf(`cdx(){
	eval "$(go-cdx $@)"
}
touch %s
touch %s
`, config.BookMarkFile, config.HistoryFile)
}
