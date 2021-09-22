# mailer
**Simple commandline mail sender**

* Repo: github.com/pepa65/mailer
* Config-less (it's not a bug, it's a feature..!)
* No-install single binary
* Defaulting to gmail's smtp server

## Build
`go build -ldflags="-s -w"; upx --brute mailer`

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
