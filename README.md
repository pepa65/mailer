[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/mailer)](https://goreportcard.com/report/github.com/pepa65/mailer)
[![GoDoc](https://godoc.org/github.com/pepa65/mailer?status.svg)](https://godoc.org/github.com/pepa65/mailer)

# mailer - Simple commandline SMTP client
* **v0.7.0**
* Repo: [github.com/pepa65/mailer](https://github.com/pepa65/mailer)
* Completely config-less, send purely from the commandline
* But parameters can also be set in `.mailer` in the current directory.
* No-install single binary
* Defaulting to gmail's smtp server
* Licence: GPLv3+

## Install
```
# gobinaries.com:
wget -qO- gobinaries.com/pepa65/mailer |sh

# Go get (If Golang is installed properly):
go install github.com/pepa65/mailer@latest

# Download:
wget -O mailer https://gitlab.com/pepa65/mailer/-/jobs/artifacts/master/raw/BINARY?job=building
# replace BINARY by: mailer, mailer_pi, mailer_osx, mailer_bsd or mailer.exe

# Go build/install:
git clone https://github.com/pepa65/mailer; cd mailer; go install

# Smaller binaries:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o mailer_pi
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o mailer_bsd
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o mailer_osx
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o mailer.exe

# More extreme shrinking:
upx mailer*

# Move them in the Go-binary path (if in your PATH):
mv mailer* ~/go/bin/

# Or move to a manually managed binaries location:
sudo mv mailer* /usr/local/bin/
```

## Usage
```
mailer v0.7.1 - Simple commandline SMTP client [repo: github.com/pepa65/mailer]
Usage:  mailer ESSENTIALS BODY [OPTIONS]
    ESSENTIALS:
        -u|--user USER            For logging in to mail server. ^1
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin.
        -t|--to EMAILS            To email(s). ^2
        -s|--subject TEXTLINE     Subject line.
    BODY is either one of:
        -m|--message TEXT         Message text.
        -F|--file FILENAME        File containing the message text.
    OPTIONS:
        -a|--attachment FILE      Filename to attach [multiple flags allowed].
        -S|--server SERVER        Mail server [default: smtp.gmail.com].
        -P|--port PORT            Port, like 25 or 465 [default: 587]. ^3
        -T|--tls                  Use SSL/TLS instead of StartTLS. ^3
        -c|--cc EMAILS            Cc email(s). ^2
        -b|--bcc EMAILS           Bcc email(s). ^2
        -r|--reply EMAILS         Reply-To email(s). ^2
        -R|--read EMAILS          Email(s) to send ReadReceipts to. ^2
        -f|--from NAME|EMAIL      The name to use with the USER's email. ^1
        -h|--help                 Only show this help text.
Notes:
    1. If USER is not an email address, '-f'/'--from' should have EMAIL!
    2. Emails can be like "you@and.me" or like "Some Name <you@and.me>" and
       can be strung together comma-separated. Mind the shell's parsing!
    3. StartTLS is the default, except when PORT is 465, then SSL/TLS is used.
    4. Commandline errors print help text and the error to stdout and return 1.
    5. Errors with sending are printed to stdout and return exitcode 2.
    6. If file '.mailer' is present in PWD, its config parameters will be used.
    7. Commandline parameters take precedence over Configfile parameters.
       In the case of attachments, both sources will be used.
```

### Config file
The file `.mailer` in the current directory can be used to set some or all parameters.
See the example file in this repo.
