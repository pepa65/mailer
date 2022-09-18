package main

import (
	"bufio"
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
)

const version = "0.5.0"

var defaultport = "587"
var defaultserver = "smtp.gmail.com"

var self string

func usage() {
	fmt.Printf(`%v v%v - Commandline mail sender (repo: github.com/pepa65/mailer)
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
        -S|--server SERVER        Mail server (default: `+defaultserver+`)
        -P|--port PORT            Port, like 25 or 465 (default: `+defaultport+`)
        -T|--tls                  Use SSL/TLS instead of (the default) StartTLS
        -c|--cc EMAILS            Cc email(s)
        -b|--bcc EMAILS           Bcc email(s)
        -r|--reply EMAILS         Reply-To email(s)
    - Emails can be like "you@and.me" or like "Some String <you@and.me>",
      and need to be comma-separated. Any argument must survive shell-parsing!
    - Commandline errors print help and the error to stdout and return 1.
    - Send errors print to stdout and return exitcode 2.
`, self, version)
}

func errormsg(msg string) {
	usage()
	fmt.Printf("\nERROR %v\n", msg)
	os.Exit(1)
}

func main() {
	// Use name of binary/link
	self = os.Args[0]
	i := strings.IndexByte(self, '/')
	for i >= 0 {
		self = self[i+1:]
		i = strings.IndexByte(self, '/')
	}

	// Usage on no arguments
	nArgs := len(os.Args)
	if nArgs == 1 {
		usage()
		return
	}

	// Parse commandline
	i = 1
	var from, to, subject, user, password, server, port, cc, bcc, reply, message, file string
	var ssltls bool
	for i < nArgs {
		switch os.Args[i] {
		case "-f", "--from":
			if from != "" {
				errormsg("Can't use -f/--from twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -f/--from must have an argument")
			}
			from = os.Args[i+1]
			i = i + 1
		case "-t", "--to":
			if to != "" {
				errormsg("Can't use -t/--to twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -t/--to must have an argument")
			}
			to = os.Args[i+1]
			i = i + 1
		case "-s", "--subject":
			if subject != "" {
				errormsg("Can't use -s/--subject twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -s/--subject must have an argument")
			}
			subject = os.Args[i+1]
			i = i + 1
		case "-u", "--user":
			if user != "" {
				errormsg("Can't use -u/--user twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -u/--user must have an argument")
			}
			user = os.Args[i+1]
			i = i + 1
		case "-p", "--password":
			if password != "" {
				errormsg("Can't use -p/--password twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -p/--password must have an argument")
			}
			password = os.Args[i+1]
			i = i + 1
		case "-S", "--server":
			if server != "" {
				errormsg("Can't use -S/--server twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -S/--server must have an argument")
			}
			server = os.Args[i+1]
			i = i + 1
		case "-P", "--port":
			if port != "" {
				errormsg("Can't use -P/--port twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -P/--port must have an argument")
			}
			port = os.Args[i+1]
			i = i + 1
		case "-T", "--tls":
			ssltls = true
		case "-c", "--cc":
			if cc != "" {
				errormsg("Can't use -c/--cc twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -c/--cc must have an argument")
			}
			cc = os.Args[i+1]
			i = i + 1
		case "-b", "--bcc":
			if bcc != "" {
				errormsg("Can't use -b/--bcc twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -b/--bcc must have an argument")
			}
			bcc = os.Args[i+1]
			i = i + 1
		case "-r", "--reply":
			if reply != "" {
				errormsg("Can't use -r/--reply twice")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -r/--reply must have an argument")
			}
			reply = os.Args[i+1]
			i = i + 1
		case "-m", "--message":
			if message != "" {
				errormsg("Can't use -m/--message twice")
			}
			if file != "" {
				errormsg("Can't use both -m/--message and -F/--file flags")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -m/--message must have an argument")
			}
			message = os.Args[i+1]
			i = i + 1
		case "-F", "--file":
			if file != "" {
				errormsg("Can't use -F/--file twice")
			}
			if message != "" {
				errormsg("can't use both -m/--message and -F/--file flags")
			}
			if i+2 > len(os.Args) {
				errormsg("Flag -F/--file must have an argument")
			}
			file = os.Args[i+1]
			i = i + 1
		case "-h", "--help":
			usage()
			return
		default:
			errormsg("unknown commandline option: " + os.Args[i])
		}
		i += 1
	}
	if to == "" {
		errormsg("Mandatory option -t/--to missing")
	}
	if subject == "" {
		errormsg("Mandatory option -s/--subject missing")
	}
	if user == "" {
		errormsg("Mandatory option -u/--user missing")
	}
	if password == "" {
		errormsg("Mandatory option -p/--password missing")
	}
	if server == "" {
		server = defaultserver
	}
	if port == "" {
		port = defaultport
	}
	if message == "" && file == "" {
		errormsg("Content missing, neither -m/--message nor -F/--file given")
	}
	if password == "-" {
		pwd := []byte{}
		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			pwd = append(pwd, in.Bytes()...)
		}
		password = string(pwd)
	}

	if from == "" { // FROM is USER if not given
		from = user
	} else if strings.Index(from, "@") < 0 { // FROM includes USER as email if no email given
		from += " <" + user + ">"
	}

	// Populate email
	if message == "" {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			errormsg("File not found: '" + file + "'")
		}
		message = string(f)
	}
	mail := email.NewEmail()
	mail.From = from
	mail.To = strings.Split(to, ",")
	if bcc != "" {
		mail.Bcc = strings.Split(bcc, ",")
	}
	if cc != "" {
		mail.Cc = strings.Split(cc, ",")
	}
	if reply != "" {
		mail.ReplyTo = strings.Split(reply, ",")
	}
	mail.Subject = subject
	mail.Text = []byte(message)

	// Send mail
	serverport := server+":"+port
	auth := smtp.PlainAuth("", user, password, server)
	tc := &tls.Config{ServerName:server, InsecureSkipVerify:false}
	var err error
	if ssltls {
		err = mail.SendWithTLS(serverport, auth, tc)
	} else {
		err = mail.SendWithStartTLS(serverport, auth, tc)
	}
	if err != nil { // No news is good news
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
}
