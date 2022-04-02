# mailer
## Simple commandline mail sender
* **v0.1.0**
* Repo: [github.com/pepa65/mailer](https://github.com/pepa65/mailer)
* Contact: pepa65 <pepa65@passchier.net>
* Config-less (it's not a bug, it's a feature..!)
* No-install single binary
* Defaulting to gmail's smtp server

## Install
```
# gobinaries.com:
wget -qO- gobinaries.com/pepa65/mailer |sh

# Go get (If Golang is installed properly):
go get github.com/pepa65/mailer

# Download:
wget https://gitlab.com/pepa65/mailer/-/jobs/artifacts/master/raw/BINARY?job=building
# replace BINARY by: mailer, mailer_pi, mailer_osx, mailer_bsd or mailer.exe

# Go build/install:
git clone https://github.com/pepa65/mailer; cd mailer; go install

# Smaller binaries:
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o mailer_pi
GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o mailer_bsd
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o mailer_osx
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o mailer.exe

# More extreme shrinking:
upx --brute mailer*

# Move them in the Go-binary path:
mv mailer* ~/go/bin/

## Usage
```
mailer v0.1.0 - Simple commandline mail sender (repo: github.com/pepa65/mailer)
Usage:  mailer CONTENT MANDATORIES [OPTIONALS]
    CONTENT is either one of:
        -m|--message TEXT         Message text
        -F|--file FILENAME        File containing the message text
    MANDATORIES:
        -t|--to EMAILS            To email(s)
        -s|--subject TEXT         Subject line
        -u|--user USER            For logging in to mail server (*)
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin
    OPTIONALS:
        -S|--server SERVER        Mail server (default: smtp.gmail.com)
        -P|--port PORT            Port, like 25 or 465 (default: 587)
        -c|--cc EMAILS            Cc email(s)
        -b|--bcc EMAILS           Bcc email(s)
        -r|--reply EMAILS         Reply-To email(s)
        -f|--from NAME            The name to use with the USER's email (*)
    (*) If USER is not an email address, NAME should contain one!
    Emails can be like "you@and.me" or like "Some String <you@and.me>",
    and need to be comma-separated. Any arguments must survive shell-parsing!
```
