package cmd

type Config struct {
	HistoryFile  string
	BookMarkFile string
	Command      string
	NoOutput     bool
	UseSSH       bool
	Make         bool
	CustomSource []CustomSource
	FuzzyFinder  FuzzyFinder
}

type FuzzyFinder struct {
	CommandPath string
	Option      string
}

type CustomSource struct {
	Name    string
	Path    string
	Command string
	Format  string
}

type Record struct {
	Number int
	Path   string
}
