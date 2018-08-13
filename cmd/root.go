// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-cdx",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// Initを出力して終了する
		if isInit {
			PrintInitText()
		}

		if popd {
			fmt.Print("popd")
			os.Exit(0)
		}

		if fuz, err := getCdCommandWithFinder(); err == nil {
			fmt.Println(fuz)
		} else {
			p, _ := os.Getwd()
			if len(args) != 0 {
				p, _ = homedir.Expand(args[0])
			}

			// Bookmarkして終了
			if addBookmark {
				AppendRecord(p, config.BookMarkFile)
				os.Exit(0)
			}

			if com, err := getCdCommand(p, os.Stderr, os.Stdin); err != nil {
				Fatal(err)
			} else {
				fmt.Println(com)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		fmt.Sprint("config file default :", filepath.Join(home, ".config", "cdx", ".go-cdx.json")))
	rootCmd.Flags().Bool("help", false, "help")

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	})

	//history
	rootCmd.Flags().BoolVarP(&useHistory, "history", "h", false, "ブックマークからcdxします")
	//bookmark
	rootCmd.Flags().BoolVarP(&useBookmark, "bookmark", "b", false, "ブックマークからcdxします")
	//make
	rootCmd.Flags().Bool("make", false, "ディレクトリが無い場合、作ってから移動します")
	viper.BindPFlag("make", rootCmd.Flags().Lookup("make"))
	//no-output
	rootCmd.Flags().Bool("no-output", false, "Stdoutに何も出力しません")
	viper.BindPFlag("NoOutput", rootCmd.Flags().Lookup("no-output"))
	//custom
	rootCmd.Flags().StringVarP(&customSource, "custom", "c", "", "コマンドの出力からcdxします")
	//add bookmark
	rootCmd.Flags().BoolVar(&addBookmark, "add", false, "カレントディレクトリをBookmarkします")
	//popd
	rootCmd.Flags().BoolVarP(&popd, "popd", "p", false, "popdを使います")
	//init
	rootCmd.Flags().BoolVar(&isInit, "init", false, "evalすることでcdxを使えるようにするコマンド列を出力します")
}

// flags
var useHistory, useBookmark, addBookmark, popd, isInit bool
var customSource string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-cdx" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config", "go-cdx"))
		viper.SetConfigName(".go-cdx")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetDefault("BinaryPath", filepath.Join(os.Getenv("GOPATH"), "bin", "go-cdx"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

}
