package cmd

import (
	"fmt"
	"os"
)

const (
	dev    = "Development"
	beta   = "Beta"
	stable = "Stable"
)

type CdxVersion struct {
	Major  int
	Minor  int
	Build  int
	Status string
	Day    string
}

var Version = CdxVersion{
	Major:  2,
	Minor:  1,
	Build:  5,
	Status: stable,
	Day:    "2019/09/16",
}

func (v CdxVersion) Print() {
	fmt.Fprintf(os.Stderr, "cdx v%d.%d.%d %s(%s)\n", v.Major, v.Minor, v.Build, v.Status, v.Day)
}
