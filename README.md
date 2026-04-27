[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/mailer)](https://goreportcard.com/report/github.com/pepa65/mailer)
<img src="https://raw.githubusercontent.com/pepa65/mailer/master/mailer.png" width="120" alt="mailer icon" align="right">

# mailer - Simple commandline SMTP client
* **v1.7.0**
* Repo: [github.com/pepa65/mailer](https://github.com/pepa65/mailer)
* No-install single binary.
* Completely config-less, can send purely from the commandline.
* But parameters can also be set in a configfile.
* Can send plaintext, html, or both, and attachments.
* Licence: GPLv3+

## Install
```
# Download (replace BINARY by: mailer, mailer_pi, mailer_osx, mailer_bsd or mailer.exe)
wget -O BINARY https://gitlab.com/pepa65/mailer/-/jobs/artifacts/master/raw/BINARY?job=building

# Go get (If Golang is installed properly)
go install github.com/pepa65/mailer@latest

# Go clone/install (If Golang is installed properly)
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
mailer v1.7.0 - Simple commandline SMTP client [repo: github.com/pepa65/mailer]
Usage:  mailer [ESSENTIALS] [BODY] [OPTIONS]
    ESSENTIALS (like any option, can be set in a configfile):
        -u|--user USER           For logging in to mail server.
        -p|--password PASSWORD   If PASSWORD is '-', it is read from stdin.
        -e|--email EMAIL         The From email. ^1
        -t|--to EMAILS           The To email(s). ^1 ^2
        -s|--subject TEXTLINE    The Subject line.
    BODY (can be both plaintext and html, but each from either string or file):
        -m|--message PLAINTEXT   Message string in plain text.
        -M|--messagefile FILE    File containing the plain text message.
        -f|--formatted HTML      Message string in html.
        -F|--formattedfile FILE  File containing the html message.
    OPTIONS:
        -o|--options CONFIGFILE  File with options [default: mailer.yml]. ^3
        -a|--attachment FILE     File to attach [multiple flags allowed]. ^4
        -S|--server SERVER       Mail server [default: smtp.gmail.com].
        -P|--port PORT           Port, like 25 or 465 [default: 587]. ^5
        -T|--tls                 Use SSL/TLS instead of StartTLS. ^5
        -c|--cc EMAILS           The Cc email(s). ^1 ^2
        -b|--bcc EMAILS          The Bcc email(s). ^1 ^2
        -r|--reply EMAILS        The Reply-To email(s). ^1 ^2
        -R|--read EMAILS         The email(s) to send ReadReceipts to. ^1 ^2
        -U|--unsubscribe         Comma-separated unsubscribe targets. ^6
        -V|--version             Only show the version.
        -h|--help                Only show this help text.
Notes:
    - Commandline options take precedence over CONFIGFILE options.
    - Commandline errors print help text and the error to stdout and return 1.
    - Sending errors are printed to stdout and return exitcode 2.
    1. EMAIL can be like "you@and.me" or like "Some String <you@and.me>".
    2. EMAILS must be comma-separated. (Easiest to enclose EMAILS in quotes.)
    3. Could be the only option, if all ESSENTIALS and BODY options get set.
       If the default CONFIGFILE exists, no Commandline arguments are needed.
    4. All FILEs given in the CONFIGFILE and on the commandline will be used.
    5. If -T/--tls is used, the default port switches to 465.
    6. Targets are email-addresses or URLs (no 'mailto:', 'https://' or '<>').
```

Note that in case of the `-U`/`--unsubscribe` flag, that besides a List-Unsubscribe header,
also a List-Unsubscribe-Post header is added if a URL is given, to make it One-Click.
This does require that the URL accepts POST, so it does not require further interaction.
The To-address will be added at the end of the given URL, so structure the end like: `?email=`,
and there should be only 1 To-address for this to be meaningful!

### Configfile
The file given after `-o`/`--options` can be used to set some or all options,
see the example file `mailer.yml` in this repo. The default is `$PWD/mailer.yml`.
The field names are the same as the long option flags.

The YAML syntax for including blocks of text is tricky, using files instead is more predictable.
When using `|+` to include blocks of text, note that `: ` (colon-space) and ` #` (space-hash)
are likely to cause a YAML syntax error... Replace space with ` ` (no-break space, U+00A0).
YAML also supports various quoting options, where a newline gets inserted on an empty line.

