package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
)

const version = "0.1.0"
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
        -t|--to EMAILS            To email(s)
        -s|--subject TEXT         Subject line
        -u|--user USER            For logging in to mail server (*)
        -p|--password PASSWORD    If PASSWORD is a dash, it is read from stdin
    OPTIONALS:
        -S|--server SERVER        Mail server (default: ` + defaultserver + `)
        -P|--port PORT            Port, like 25 or 465 (default: ` + defaultport + `)
        -c|--cc EMAILS            Cc email(s)
        -b|--bcc EMAILS           Bcc email(s)
        -r|--reply EMAILS         Reply-To email(s)
        -f|--from NAME            The name to use with the USER's email (*)
    (*) If USER is not an email address, NAME should contain one!
    Emails can be like "you@and.me" or like "Some String <you@and.me>",
    and need to be comma-separated. Any arguments must survive shell-parsing!
`, self, version);
}

func error(msg string) {
	usage()
	fmt.Printf("\nERROR %v\n", msg)
	os.Exit(1)
}

func extract(emails string) []string {
	toall := strings.Split(emails, ",")
	for i,email := range toall {
		a := strings.Index(email, "<")
		b := strings.Index(email, ">")
		if a + 1 < b && a + 1 > 0  {
			toall[i] = email[a + 1:b]
		}
	}
	return toall
}

func main() {
	// Use name of binary/link
  self = os.Args[0]
  i := strings.IndexByte(self, '/')
  for i >= 0 {
    self = self[i + 1:]
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
	from,to,subject,user,password,server,port,cc,bcc,reply,message,file := "","","","","","","","","","","",""
	for i < nArgs {
		switch os.Args[i] {
		case "-f","--from":
			if from != "" {
				error("Can't use -f/--from twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -f/--from must have an argument")
			}
			from = os.Args[i + 1]
			i = i + 1
		case "-t","--to":
			if to != "" {
				error("Can't use -t/--to twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -t/--to must have an argument")
			}
			to = os.Args[i + 1]
			i = i + 1
		case "-s","--subject":
			if subject != "" {
				error("Can't use -s/--subject twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -s/--subject must have an argument")
			}
			subject = os.Args[i + 1]
			i = i + 1
		case "-u","--user":
			if user != "" {
				error("Can't use -u/--user twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -u/--user must have an argument")
			}
			user = os.Args[i + 1]
			i = i + 1
		case "-p","--password":
			if password != "" {
				error("Can't use -p/--password twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -p/--password must have an argument")
			}
			password = os.Args[i + 1]
			i = i + 1
		case "-S","--server":
			if server != "" {
				error("Can't use -S/--server twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -S/--server must have an argument")
			}
			server = os.Args[i + 1]
			i = i + 1
		case "-P","--port":
			if port != "" {
				error("Can't use -P/--port twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -P/--port must have an argument")
			}
			port = os.Args[i + 1]
			i = i + 1
		case "-c","--cc":
			if cc != "" {
				error("Can't use -c/--cc twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -c/--cc must have an argument")
			}
			cc = os.Args[i + 1]
			i = i + 1
		case "-b","--bcc":
			if bcc != "" {
				error("Can't use -b/--bcc twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -b/--bcc must have an argument")
			}
			bcc = os.Args[i + 1]
			i = i + 1
		case "-r","--reply":
			if reply != "" {
				error("Can't use -r/--reply twice")
			}
			if i + 2 > len(os.Args) {
				error("Flag -r/--reply must have an argument")
			}
			reply = os.Args[i + 1]
			i = i + 1
		case "-m","--message":
			if message != "" {
				error("Can't use -m/--message twice")
			}
			if file != "" {
				error("Can't use both -m/--message and -F/--file flags")
			}
			if i + 2 > len(os.Args) {
				error("Flag -m/--message must have an argument")
			}
			message = os.Args[i + 1]
			i = i + 1
		case "-F","--file":
			if file != "" {
				error("Can't use -F/--file twice")
			}
			if message != "" {
				error("can't use both -m/--message and -F/--file flags")
			}
			if i + 2 > len(os.Args) {
				error("Flag -F/--file must have an argument")
			}
			file = os.Args[i + 1]
			i = i + 1
		case "-h","--help":
			usage()
			return
		default:
			error("unknown commandline option: " + os.Args[i])
		}
		i += 1
	}
	if to == "" {
		error("Mandatory option -t/--to missing")
	}
	if subject == "" {
		error("Mandatory option -s/--subject missing")
	}
	if user == "" {
		error("Mandatory option -u/--user missing")
	}
	if password == "" {
		error("Mandatory option -p/--password missing")
	}
	if message == "" && file == "" {
		error("Content missing, neither -m/--message nor -F/--file given")
	}
	if server == "" {
		server = defaultserver
	}
	if port == "" {
		port = defaultport
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
	// Assemble complete To list of emails
	all := to
	if cc != "" {
		all += "," + cc
	}
	if bcc != "" {
		all += "," + bcc
	}
	toall := extract(all)

	// Assemble email with headers and body
	body := "To: " + to + "\nFrom: " + from + "\n"
	if reply !="" {
		body = body + "Reply-To: " + reply + "\n"
	}
	if cc !="" {
		body += "Cc: " + cc + "\n"
	}
	body += "Subject: " + subject + "\n\n"
	if message != "" {
		body += message
	} else {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			error("File not found: '" + file + "'")
		}
		body += string(f)
	}

	// Send with net/smtp.sendmail
	auth := smtp.PlainAuth("", user, password, server)
	err := smtp.SendMail(server + ":" + port, auth, user, toall, []byte(body))
	if err != nil { // No news is good news
		fmt.Println("Error: ", err)
	}
}
