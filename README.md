[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/mailer)](https://goreportcard.com/report/github.com/pepa65/mailer)
[![GoDoc](https://godoc.org/github.com/pepa65/mailer?status.svg)](https://godoc.org/github.com/pepa65/mailer)

# mailer - Simple commandline SMTP client
* **v0.6.0**
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
mailer v0.6.0 - Simple commandline SMTP client (repo: github.com/pepa65/mailer)
Usage:  mailer CONTENT MANDATORIES [OPTIONALS]
    CONTENT is either one of:
        -m|--message TEXT         Message text.
        -F|--file FILENAME        File containing the message text.
    MANDATORIES:
        -t|--to EMAILS            To email(s). [1]
        -s|--subject TEXT         Subject line.
        -u|--user USER            For logging in to mail server.
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin.
    OPTIONALS:
        -S|--server SERVER        Mail server (default: smtp.gmail.com).
        -P|--port PORT            Port, like 25 or 465 (default: 587). [3]
        -T|--tls                  Use SSL/TLS instead of StartTLS. [3]
        -c|--cc EMAILS            Cc email(s).
        -b|--bcc EMAILS           Bcc email(s).
        -r|--reply EMAILS         Reply-To email(s).
        -f|--from NAME            The name to use with the USER's email.
Notes:
    1. If USER is not an email address, NAME should contain one!
    2. Emails can be like "you@and.me" or like "Some String <you@and.me>" and
       need to be comma-separated. Any argument must survive shell-parsing!
    3. StartTLS is the default, except when PORT is 465, then SSL/TLS is used.
    4. Commandline errors print help text and the error to stdout and return 1.
    5. Errors with sending are printed to stdout and return exitcode 2.
    6. If ".mailer" is present in PWD it will be read for config parameters.
```

### Config file
The file `.mailer` in the current directory can be used to set some or all parameters.
See the example file in this repo.
