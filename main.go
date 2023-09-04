// main.go - Simple commandline SMTP client

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v2"
)

const version = "0.6.0"

type Config struct {
	User     string
	Password string
	Server   string
	Port     string
	TLS      string
	From     string
	CC       string
	BCC      string
	Reply    string
	To       string
	Subject  string
	Message  string
	File     string
}

var defaultport = "587"
var defaultserver = "smtp.gmail.com"
var self string

func usage() {
	fmt.Printf(`%v v%v - Simple commandline SMTP client (repo: github.com/pepa65/mailer)
Usage:  mailer CONTENT MANDATORIES [OPTIONALS]
    CONTENT is either one of:
        -m|--message TEXT         Message text.
        -F|--file FILENAME        File containing the message text.
    MANDATORIES:
        -t|--to EMAILS            To email(s).
        -s|--subject TEXT         Subject line.
        -u|--user USER            For logging in to mail server. [1]
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin.
    OPTIONALS:
        -S|--server SERVER        Mail server (default: `+defaultserver+`).
        -P|--port PORT            Port, like 25 or 465 (default: `+defaultport+`). [3]
        -T|--tls                  Use SSL/TLS instead of StartTLS. [3]
        -c|--cc EMAILS            Cc email(s).
        -b|--bcc EMAILS           Bcc email(s).
        -r|--reply EMAILS         Reply-To email(s).
        -f|--from NAME            The name to use with the USER's email.
        -h|--help                 Only show this help text.
Notes:
    1. If USER is not an email address, NAME should contain one!
    2. Emails can be like "you@and.me" or like "Some String <you@and.me>" and
       need to be comma-separated. Any arguments must survive shell-parsing!
    3. StartTLS is the default, except when PORT is 465, then SSL/TLS is used.
    4. Commandline errors print help text and the error to stdout and return 1.
    5. Errors with sending are printed to stdout and return exitcode 2.
    6. If file ".mailer" is present in PWD its config parameters will be used.
    7. Commandline parameters take precedence over Configfile parameters.
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

	// Use .mailer if present
	nArgs := len(os.Args)
	var cfg Config
	_, err := os.Stat(".mailer")
	if err == nil { // .mailer present
		cfgdata, err := ioutil.ReadFile(".mailer")
		if err != nil {
			errormsg("Config file '.mailer' not found")
		}
		err = yaml.UnmarshalStrict(cfgdata, &cfg)
		if err != nil {
			errormsg("Error in config file '.mailer':\n"+err.Error())
		}
	} else { // No .mailer in PWD
		if nArgs == 1 { // Usage on no arguments
			usage()
			return
		}
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
	if cc == "" {
		cc = cfg.CC
	}
	if bcc == "" {
		bcc = cfg.BCC
	}
	if reply == "" {
		reply = cfg.Reply
	}
	if to == "" {
		to = cfg.To
	}
	if to == "" {
		errormsg("Mandatory option 'to' missing")
	}
	if subject == "" {
		subject = cfg.Subject
	}
	if subject == "" {
		errormsg("Mandatory option 'subject' missing")
	}
	if user == "" {
		user = cfg.User
	}
	if user == "" {
		errormsg("Mandatory option 'user' missing")
	}
	if password == "" {
		password = cfg.Password
	}
	if password == "" {
		errormsg("Mandatory option 'password' missing")
	}
	if server == "" {
		server = cfg.Server
	}
	if server == "" {
		server = defaultserver
	}
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = defaultport
	}
	if port == "465" {
		ssltls = true
	}
	if message == "" && file == "" && cfg.Message != "" && cfg.File != "" {
		errormsg("Can't have both 'message' and 'file' options set in .mailer")
	}
	if message == "" {
		message = cfg.Message
	}
	if file == "" {
		file = cfg.File
	}
	if message == "" && file == "" {
		errormsg("Content missing, neither 'message' nor 'file' option given")
	}
	if message != "" && file != "" { // FILE (from commandline) takes precedence
		message = ""
	}
	if password == "-" {
		pwd := []byte{}
		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			pwd = append(pwd, in.Bytes()...)
		}
		password = string(pwd)
	}

	if from == "" {
		from = cfg.From
	}
	if !strings.Contains(from, "@") { // FROM does not contains email
		if from == "" { // FROM not given, must use USER
			from = user
		} else { // FROM given as a NAME
			from += " <" + user + ">"
		}
		if !strings.Contains(from, "@") { // No FROM email, also not in USER
			errormsg("No 'from' email nor 'user' email")
		}
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
	serverport := server + ":" + port
	auth := smtp.PlainAuth("", user, password, server)
	tc := &tls.Config{ServerName: server, InsecureSkipVerify: false}
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
