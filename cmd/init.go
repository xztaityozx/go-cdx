package cmd

import (
	"fmt"
	"runtime"

	"github.com/xztaityozx/go-cdx/environment"
)

func PrintInit(env environment.Environment) {
	fmt.Printf(func() string {
		if runtime.GOOS == "windows" {
			return fmt.Sprint(`function cdx(){
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
`)
		} else {
			return fmt.Sprint(`function cdx(){
eval "$(go-cdx $@)"
}
[ -f %s ] || touch %s
[ -f %s ] || touch %s
[ -d $HOME/.config/go-cdx ] || mkdir -p $HOME/.config/go-cdx`)
		}
	}(), cfg.File.History, cfg.File.History, cfg.File.BookMark, cfg.File.BookMark)
}
