# goxion

## Description

某VPNに簡単にログインできるようになります。

## Installation

### For Mac

selenium-server-standaloneをまずインストールしてください。

```bash
$ brew install selenium-server-standalone
```

次に、SafariDriverを以下のページから取得し、インストールしてください。

* http://www.seleniumhq.org/download/

#### Safariの環境設定

* 環境設定 - 詳細 から 'メニューバーに"開発"メニューを表示'にチェック
* 開発 - リモートオートメーションを許可

### For Windows

お使いのPCに合ったIE用のWebDriverを以下のページから取得し、インストールしてください。

* http://www.seleniumhq.org/download/

## Usage

```bash
$ goxion -h                                                               ⏎
NAME:
   goxion

USAGE:
   goxion [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   tenten0213

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value    userId for login.
   --password value, -p value    password for login.
   --pjid value, -i value        projectid for login.
   --url value, -r value         login page's url.
   --matrixpath value, -m value  path to matrix file.
   --help, -h                    show help
   --version, -v                 print the version
```

以下のように実行できます。

`$ goxion -u username -p password -i pjid -r http://vpn.example.com -m ./matrix.csv`

また、各パラメータは以下の環境変数に設定することが出来ます。
環境変数を設定した場合は、単に`goxion`と実行することが可能です。

* `GOXION_USER`
* `GOXION_PASSWORD`
* `GOXION_URL`
* `GOXION_MATRIX`

## Build

```bash
$ go get github.com/sclevine/agouti
$ go get github.com/urfave/cli
$ go get github.com/mitchellh/go-ps
```

### For Mac

```bash
$ GOOS=darwin GOARCH=amd64 go build -o goxion main.go
```

### For Windows

```bash
# for 64bit
$ GOOS=windows GOARCH=amd64 go build -o goxion.exe main.go
# for 32bit
$ GOOS=windows GOARCH=386 go build -o goxion.exe main.go
```

## Author

[tenten0213](https://github.com/tenten0213)

