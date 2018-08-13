package cmd

import (
	"fmt"
	"os"
)

var version Version = Version{
	Major:  1,
	Minor:  0,
	Build:  10,
	Status: "Beta",
	Date:   "2018/08/13",
}

type Version struct {
	Major  int
	Minor  int
	Build  int
	Status string
	Date   string
}

func (v Version) ToString() string {
	return fmt.Sprintf("cdx version %d.%d.%d %s (%s)\n\nAuthor: xztaityozx\nRepository: https://github.com/xztaityozx/go-cdx\n\nLicense: MIT", v.Major, v.Minor, v.Build, v.Status, v.Date)
}

func PrintVersion() {
	fmt.Print(version.ToString())
	os.Exit(0)
}
