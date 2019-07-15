// Copyright © 2019 xztaityozx
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
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/go-cdx/cd"
	"github.com/xztaityozx/go-cdx/config"
	"github.com/xztaityozx/go-cdx/customsource"
	"github.com/xztaityozx/go-cdx/environment"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-cdx",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		env := environment.NewEnvironment()

		// init
		if init, _ := cmd.Flags().GetBool("init"); init {
			PrintInit(env)
			return
		}
		// version
		if v, _ := cmd.Flags().GetBool("version"); v {
			Version.Print()
			os.Exit(1)
		}
		// add
		if add, _ := cmd.Flags().GetBool("add"); add {
			wd, _ := os.Getwd()
			cfg.File.AppendBookmark(wd)
			return
		}
		// popd
		if popd, _ := cmd.Flags().GetBool("popd"); popd {
			fmt.Print("popd")
			return
		}

		// AddHistorys
		custom, _ := cmd.Flags().GetString("custom")
		if h, _ := cmd.Flags().GetBool("history"); h {
			custom += "h"
			cfg.CustomSources = append(cfg.CustomSources,
				customsource.CustomSource{
					Name:        "history",
					SubName:     "h",
					Command:     "/bin/cat " + cfg.File.History,
					BeginColumn: 0,
				})
		}
		if b, _ := cmd.Flags().GetBool("bookmark"); b {
			custom += "b"
			cfg.CustomSources = append(cfg.CustomSources,
				customsource.CustomSource{
					Name:        "bookmark",
					SubName:     "b",
					Command:     "/bin/cat " + cfg.File.BookMark,
					BeginColumn: 0,
				})
		}

		// ここからFuzzyFinderを使う

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGSTOP)
		defer close(ch)
		go func() {
			<-ch
			cancel()
		}()
		var path string
		if len(custom) == 0 {
			// 引数のパスを使う
			path, _ = homedir.Expand(strings.Join(args, " "))
		} else {

			logrus.Info(cfg.CustomSources[0])

			// CustomSourceを使う
			source, err := customsource.New(custom, cfg.CustomSources)
			if err != nil {
				logrus.WithError(err).Error("failed search CustomSource")
				return
			}

			path, err = cfg.FuzzyFinder.Run(ctx, source, env)
			if err != nil {
				logrus.Error(err)
			}
		}

		command := cd.New(cfg.Command, path, cfg.NoOutput, cfg.Make, cfg.FuzzyFinder)
		fmt.Print(command.BuildCommand(ctx, env))

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/go-cdx/go-cdx.yml)")

	rootCmd.Flags().Bool("help", false, "help")

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	})

	// CustomSource
	rootCmd.Flags().StringP("custom", "c", "", "CustomSourceからcdします")

	// NoOutput
	rootCmd.Flags().Bool("no-output", false, "STDOUTに何も出力しません")
	viper.BindPFlag("NoOutput", rootCmd.Flags().Lookup("no-output"))

	// history
	rootCmd.Flags().BoolP("history", "h", false, "履歴からcdします")
	// bookmark
	rootCmd.Flags().BoolP("bookmark", "b", false, "ブックマークからcdします")

	// popd
	rootCmd.Flags().BoolP("popd", "p", false, "popdします")

	// add bookmark
	rootCmd.Flags().Bool("add", false, "bookmarkにカレントディレクトリを追加します")

	// init
	rootCmd.Flags().Bool("init", false, "evalすることでcdxを使えるようにするコマンド列を出力します")

	// version
	rootCmd.Flags().BoolP("version", "v", false, "versionを出力して終了します")

	// make
	rootCmd.Flags().Bool("make", false, "ディレクトリが無い場合、作ってから移動します")
	viper.BindPFlag("make", rootCmd.Flags().Lookup("make"))

}

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
		viper.SetConfigName("go-cdx")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatal("[cdx] failed read config file")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.WithError(err).Fatal("[cdx] failed unmarshal config file")
	}
}
