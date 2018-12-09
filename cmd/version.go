package cmd

import (
	"fmt"
	"os"
)

var version Version = Version{
	Major:  2,
	Minor:  0,
	Build:  25,
	Status: "Beta",
	Date:   "2018/12/09",
}

type Version struct {
	Major  int
	Minor  int
	Build  int
	Status string
	Date   string
}

func (v Version) ToString() string {
	return fmt.Sprintf(`cdx version %d.%d.%d %s (%s)
	
Author:     xztaityozx
Repository: https://github.com/xztaityozx/go-cdx
License:    MIT`, v.Major, v.Minor, v.Build, v.Status, v.Date)
}

func PrintVersion() {
	os.Stderr.WriteString(version.ToString())
	os.Stderr.Close()
	os.Exit(0)
}
