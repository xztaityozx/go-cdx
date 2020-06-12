package config

import (
	"fmt"
	"runtime"
)

func Initialize() error {

	if runtime.GOOS == "windows" {
		fmt.Println(`
function cdx() {
    begin{
        [System.Collections.Generic.List[string]]$paths=@{}
    }
    process{
        $paths.Add($_)
    }
    end{
        $command = "$($paths | go-cdx $args)";
        Invoke-Expression "$command"
    }
}`)
	} else {
		fmt.Println(`function() cdx() {command="$(go-cdx $@)"; eval "${command}"}`)
	}

	return nil
}
