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
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/go-cdx/cd"
	"github.com/xztaityozx/go-cdx/config"
	"github.com/xztaityozx/go-cdx/subcmd"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg config.Config

var rootCmd = &cobra.Command{
	Use:     "go-cdx",
	Short:   "",
	Long:    ``,
	Version: "2.2.1",
	PreRun:  subcmd.GenCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range []struct {
			name   string
			action func() (string, error)
		}{
			{name: "add", action: func() (string, error) {
				return subcmd.Add(cfg.BookmarkFile)
			}},
			{name: "popd", action: subcmd.Popd},
			{name: "init", action: subcmd.Initialize},
			{name: "git-root", action: subcmd.GitRoot},
		} {
			if f, _ := cmd.Flags().GetBool(v.name); f {
				if command, err := v.action(); err != nil {
					logrus.WithError(err).Fatal("[cdx] failed sub command")
				} else {
					fmt.Println(command)
				}
				return
			}
		}

		cs, _ := cmd.Flags().GetString("cdxsource")
		if f, _ := cmd.Flags().GetBool("history"); f {
			cs += "h"
		}
		if f, _ := cmd.Flags().GetBool("bookmark"); f {
			cs += "b"
		}

		// list up candidate paths
		home, _ := homedir.Dir()
		var can []string
		for _, v := range args {
			can = append(can, strings.Replace(v, "~", home, 1))
		}
		if f, _ := cmd.Flags().GetBool("stdin"); f {
			scan := bufio.NewScanner(os.Stdin)
			for scan.Scan() {
				can = append(can, scan.Text())
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		defer close(sigCh)
		signal.Notify(sigCh, syscall.SIGINT)
		go func() {
			<-sigCh
			cancel()
		}()

		// output command string
		command, err := cd.New(cfg, can).Build(ctx, cs)
		if err != nil {
			logrus.WithError(err).Fatal(err)
			fmt.Println(config.ExitCommand())
		}

		fmt.Println(command)
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
		_ = cmd.Usage()
		os.Exit(1)
	})
	// CustomSource
	rootCmd.Flags().StringP("cdxsource", "c", "", "CustomSourceからcdします")
	// NoOutput
	rootCmd.Flags().Bool("no-output", false, "STDOUTに何も出力しません")
	_ = viper.BindPFlag("NoOutput", rootCmd.Flags().Lookup("no-output"))
	// history
	rootCmd.Flags().BoolP("history", "h", false, "履歴からcdします")
	// bookmark
	rootCmd.Flags().BoolP("bookmark", "b", false, "ブックマークからcdします")
	// popd
	rootCmd.Flags().BoolP("popd", "p", false, "popdします")
	// add bookmark
	rootCmd.Flags().Bool("add", false, "bookmarkにカレントディレクトリを追加します")
	// make
	rootCmd.Flags().Bool("make", false, "ディレクトリが無い場合、作ってから移動します")
	_ = viper.BindPFlag("Make", rootCmd.Flags().Lookup("make"))
	// stdin
	rootCmd.Flags().BoolP("stdin", "i", false, "stdinから候補を受け取ります")
	// init
	rootCmd.Flags().Bool("init", false, "init用のScriptを出力します")
	// git-root
	rootCmd.Flags().BoolP("git-root", "R", false, "gitのルートディレクトリまで移動します")
	// generate completion script
	rootCmd.Flags().String("completion", "", "generate completion script for bash, zsh, fish, PowerShell")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		logrus.Fatal(err)
	}
}
