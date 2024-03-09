[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/mailer)](https://goreportcard.com/report/github.com/pepa65/mailer)
<img src="https://raw.githubusercontent.com/pepa65/mailer/master/mailer.png" width="120" alt="mailer icon" align="right">

# mailer - Simple commandline SMTP client
* **v1.0.2**
* Repo: [github.com/pepa65/mailer](https://github.com/pepa65/mailer)
* Completely config-less, send purely from the commandline.
* But parameters can also be set in `.mailer` in the current directory.
* No-install single binary.
* Defaulting to gmail's smtp server.
* Can send plaintext, html or both, and attachments.
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
mailer v1.0.2 - Simple commandline SMTP client [repo: github.com/pepa65/mailer]
Usage:  mailer ESSENTIALS BODY [OPTIONS]
    ESSENTIALS:
        -u|--user USER            For logging in to mail server. ^1
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin.
        -t|--to EMAILS            To email(s). ^2
        -s|--subject TEXTLINE     Subject line.
    BODY can be both plaintext and html, but each from either string or file:
        -m|--message PLAINTEXT    Message plain text.
        -M|--mfile FILENAME       File containing the plain text message.
        -n|--nmessage HTML        Message html.
        -N|--nfile FILENAME       File containing the html message.
    OPTIONS:
        -a|--attachment FILE      File to attach [multiple flags allowed]. ^7
        -S|--server SERVER        Mail server [default: smtp.gmail.com].
        -P|--port PORT            Port, like 25 or 465 [default: 587]. ^3
        -T|--tls                  Use SSL/TLS instead of StartTLS. ^3
        -c|--cc EMAILS            Cc email(s). ^2
        -b|--bcc EMAILS           Bcc email(s). ^2
        -r|--reply EMAILS         Reply-To email(s). ^2
        -R|--read EMAILS          Email(s) to send ReadReceipts to. ^2
        -f|--from NAME|EMAIL      The name to use with the USER's email. ^1
        -H|--header               No empty line before BODY, extended header.
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
       In case of '-a'/'--attachment'/'attachment:', both sources will be used.
```

### Config file
The file `.mailer` in the current directory can be used to set some or all parameters.
See the example file in this repo. The fields are the same as the long option flags.
The YAML syntax for including long texts is tricky, so passing files is recommended.
When using `|+` to include blocks of text, note that `: ` (colon-space) and ` #` (space-hash)
are likely to cause a YAML syntax error... Replace space with ` ` (no-break space, U+00A0).
