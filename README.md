# mailer
**Simple commandline mail sender**

* **v0.1.0**
* Repo: github.com/pepa65/mailer
* Config-less (it's not a bug, it's a feature..!)
* No-install single binary
* Defaulting to gmail's smtp server

## Install
* **gobinaries.com**: `wget -qO- gobinaries.com/pepa65/mailer |sh`
* **Go get** If [Golang](https://golang.org/) is installed properly:
  `go get github.com/pepa65/mailer`
* **Go build/install**
  - `git clone https://github.com/pepa65/mailer; cd mailer; go install`
  - Smaller binary: `go build -ldflags="-s -w"; upx --brute mailer; mv mailer ~/go/bin/`
* **Build for other architectures**
  - `GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o mailer_pi`
  - `GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o mailer_bsd`
  - `GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o mailer_osx`
  - `GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o mailer.exe`

## Usage
```
mailer - Simple commandline mail sender (repo: github.com/pepa65/mailer)
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
