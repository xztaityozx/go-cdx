# cdx
![Go](https://github.com/xztaityozx/go-cdx/workflows/Go/badge.svg?branch=master)

_cdx_ wrapper for cd(pushd) command

## Install
```sh
$ go get -u github.com/xztaityozx/go-cdx
```

or Download from GitHub releases page

## Usage

```
Usage:
  go-cdx [flags] [path]

Flags:
      --add                bookmarkにカレントディレクトリを追加します
  -b, --bookmark           ブックマークからcdします
  -c, --cdxsource string   CustomSourceからcdします
      --config string      config file (default is $HOME/.config/go-cdx/go-cdx.yml)
      --help               help
  -h, --history            履歴からcdします
      --init               init用のScriptを出力します
      --make               ディレクトリが無い場合、作ってから移動します
      --no-output          STDOUTに何も出力しません
  -p, --popd               popdします
  -i, --stdin              stdinから候補を受け取ります
  -v, --version            versionを出力して終了します
```

## Config

- Path: `$HOME/.config/go-cdx/go-cdx.yaml or json`

```yaml
historyFile: ~/.config/go-cdx/history
bookmarkFile: ~/.config/go-cdx/bookmark
noOutput: true
make: false
source:
  - name: "ghq list"
    alias: "g"
    skipcolumn: 1
    command: "ghq list|awk -F/ -v X=\"$(ghq root)\" '{print $NF, X\"/\"$0}'|column -t"
fuzzyfinder:
  command: "fzf"
```
