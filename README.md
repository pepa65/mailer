[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/mailer)](https://goreportcard.com/report/github.com/pepa65/mailer)
<img src="https://raw.githubusercontent.com/pepa65/mailer/master/mailer.png" width="120" alt="mailer icon" align="right">

# mailer - Simple commandline SMTP client
* **v1.2.1**
* Repo: [github.com/pepa65/mailer](https://github.com/pepa65/mailer)
* No-install single binary.
* Completely config-less, can send purely from the commandline.
* But parameters can also be set in a configfile.
* Can send plaintext, html, or both, and attachments.
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
upx --best --lzma mailer*

# Move them to the local binary directory (if in your PATH):
mv mailer* ~/bin/

# Or move to a manually managed binaries location:
sudo mv mailer* /usr/local/bin/
```

## Usage
```
mailer v1.2.1 - Simple commandline SMTP client [repo: github.com/pepa65/mailer]
Usage:  mailer [ESSENTIALS] [BODY] [OPTIONS]
    ESSENTIALS (like any option, can be set in a configfile):
        -u|--user USER             For logging in to mail server. ^1
        -p|--password PASSWORD     If PASSWORD is '-', it is read from stdin.
        -t|--to EMAILS             To email(s). ^2
        -s|--subject TEXTLINE      Subject line.
    BODY (can be both plaintext and html, but each from either string or file):
        -m|--message PLAINTEXT     Message string in plain text.
        -M|--mfile FILENAME        File containing the plain text message.
        -n|--nmessage HTML         Message string in html.
        -N|--nfile FILENAME        File containing the html message.
    OPTIONS:
        -o|--options CONFIGFILE    File with options [default: mailer.cfg]. ^3
        -a|--attachment FILE       File to attach [multiple flags allowed]. ^4
        -S|--server SERVER         Mail server [default: smtp.gmail.com].
        -P|--port PORT             Port, like 25 or 465 [default: 587]. ^5
        -T|--tls                   Use SSL/TLS instead of StartTLS. ^5
        -c|--cc EMAILS             Cc email(s). ^2
        -b|--bcc EMAILS            Bcc email(s). ^2
        -r|--reply EMAILS          Reply-To email(s). ^2
        -R|--read EMAILS           Email(s) to send ReadReceipts to. ^2
        -f|--from NAME|EMAIL       The name to use with the USER's email. ^1
        -h|--help                  Only show this help text.
Notes:
    - Commandline options take precedence over CONFIGFILE options.
    - Commandline errors print help text and the error to stdout and return 1.
      Errors with sending are printed to stdout and return exitcode 2.
    1. If USER is not an email address, '-f'/'--from' should have EMAIL!
    2. EMAILs can be like "you@and.me" or like "Some String <you@and.me>" and
       can be strung together comma-separated. (Mind the shell's parsing!)
    3. Could be the only option, if all ESSENTIALS and BODY options get set.
       or if the default CONFIGFILE exists, no Commandline options are needed.
    4. All given in the CONFIGFILE and on the commandline will be used.
    5. StartTLS is the default, except when PORT is 465, then SSL/TLS is used.
```

### Configfile
The file given after `-o`/`--options` can be used to set some or all options,
see the example file `mailer.cfg` in this repo.
The field names are the same as the long option flags.

The YAML syntax for including blocks of text is tricky, using files instead is more predictable.
When using `|+` to include blocks of text, note that `: ` (colon-space) and ` #` (space-hash)
are likely to cause a YAML syntax error... Replace space with ` ` (no-break space, U+00A0).
YAML also supports various quoting options, where a newline gets inserted on an empty line.

