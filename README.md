# mailer
## Simple commandline mail sender
* **v0.5.0**
* Repo: [github.com/pepa65/mailer](https://github.com/pepa65/mailer)
* Contact: pepa65 <pepa65@passchier.net>
* Config-less (it's not a bug, it's a feature..!)
* No-install single binary
* Defaulting to gmail's smtp server
* Licence: GPLv3+

## Install
```
# gobinaries.com:
wget -qO- gobinaries.com/pepa65/mailer |sh

# Go get (If Golang is installed properly):
go get github.com/pepa65/mailer

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
upx --brute mailer*

# Move them in the Go-binary path:
mv mailer* ~/go/bin/
```

## Usage
```
mailer v0.5.0 - Simple commandline mail sender (repo: github.com/pepa65/mailer)
Usage:  mailer CONTENT MANDATORIES [OPTIONALS]
    CONTENT is either one of:
        -m|--message TEXT         Message text
        -F|--file FILENAME        File containing the message text
    MANDATORIES:
        -f|--from EMAIL           From email
        -t|--to EMAILS            To email(s)
        -s|--subject TEXT         Subject line
        -u|--user USER            For logging in to mail server (*)
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin
    OPTIONALS:
        -S|--server SERVER        Mail server (default: smtp.gmail.com)
        -P|--port PORT            Port, like 25 or 465 (default: 587)
        -T|--tls                  Use SSL/TLS instead of (the default) StartTLS
        -c|--cc EMAILS            Cc email(s)
        -b|--bcc EMAILS           Bcc email(s)
        -r|--reply EMAILS         Reply-To email(s)
    - Emails can be like "you@and.me" or like "Some String <you@and.me>",
      and need to be comma-separated. Any argument must survive shell-parsing!
    - Commandline errors print help and the error to stdout and return 1.
    - Send errors print to stdout and return exitcode 2.
```
