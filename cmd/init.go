package cmd

import (
	"fmt"
	"runtime"
)

func PrintInit() {
	fmt.Printf(func() string {
		if runtime.GOOS == "windows" {
			return `function cdx(){
$command=go-cdx "$args"
	iex "$command"
}
if (!(Test-Path %s)) {
	New-Item %s
}
if (!(Test-Path %s)) {
	New-Item %s
}
if (!(Test-Path $HOME\.config\go-cdx)) {
	mkdir	$HOME\.config\go-cdx
}
`
		} else {
			return `function cdx(){
eval "$(go-cdx $@)"
}
[ -f %s ] || touch %s
[ -f %s ] || touch %s
[ -d $HOME/.config/go-cdx ] || mkdir -p $HOME/.config/go-cdx`
		}
	}(), cfg.HistoryFile, cfg.HistoryFile, cfg.BookmarkFile, cfg.BookmarkFile)
}
