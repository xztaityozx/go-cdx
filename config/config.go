package config

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Make         bool        `yaml:"make"`
		NoOutput     bool        `yaml:"noOutput"`
		Source       []CdxSource `yaml:"source"`
		HistoryFile  string      `yaml:"history"`
		BookmarkFile string      `yaml:"bookmark"`
		FuzzyFinder  FuzzyFinder `yaml:"fuzzyfinder"`
	}
)

var devNull = ""

func DevNull() string {
	if len(devNull) == 0 {
		if runtime.GOOS == "windows" {
			devNull = " > $null"
		} else {
			devNull = " > /dev/null"
		}
	}
	return devNull
}

var shell = ""

func DefaultShell() string {
	if len(shell) == 0 {
		if runtime.GOOS == "windows" {
			shell = "powershell.exe"
		} else {
			shell = "/bin/sh"
		}
	}
	return shell
}

var exitCommand = ""

func ExitCommand() string {
	if len(exitCommand) == 0 {
		if runtime.GOOS == "windows" {
			exitCommand = "throw '[cdx] failed'"
		} else {
			exitCommand = "return 1"
		}
	}
	return exitCommand
}

func Load(path string) (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return Config{}, err
	}
	if len(path) != 0 {
		viper.SetConfigFile(path)
	} else {

		// linux/macOSは$HOME/.config/go-cdx/以下を見る
		viper.AddConfigPath(filepath.Join(home, ".config", "go-cdx"))
		// Windowsなら追加で $HOME\AppData\Roaming\go-cdxも見る
		if runtime.GOOS == "windows" {
			viper.AddConfigPath(filepath.Join(home, "AppData", "Roaming", "go-cdx"))
		}
		// ファイル名は go-cdx.{json,toml,yaml}など。viperが解釈できればなんでも
		viper.SetConfigName("go-cdx")
	}

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)

	cfg.HistoryFile = strings.Replace(cfg.HistoryFile, "~", home, 1)
	cfg.BookmarkFile = strings.Replace(cfg.BookmarkFile, "~", home, 1)

	return cfg, err
}
