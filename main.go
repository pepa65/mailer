// main.go - Simple commandline SMTP client

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v2"
)

const version = "1.2.0"

type Config struct { // Options in CONFIGFILE
	User       string
	Password   string
	Server     string
	Port       string
	TLS        string
	From       string
	CC         string
	BCC        string
	Reply      string
	Read       string
	To         string
	Subject    string
	Message    string
	Mfile      string
	Nmessage   string
	Nfile      string
	Attachment string
	CSVfile    string
}

type Record map[string]string

var (
	csvheader     []string
	records       []Record
	defaultport   = "587"
	defaultserver = "smtp.gmail.com"
	self          = ""
)

func usage() {
	fmt.Printf(`%v v%v - Simple commandline SMTP client [repo: github.com/pepa65/mailer]
Usage:  mailer [ESSENTIALS BODY OPTIONS]
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
        -o|--options CONFIGFILE    File with options. ^3
        -a|--attachment FILE       File to attach [multiple flags allowed]. ^4
        -C|--csv                   Use CSV file for bulk sending.
        -S|--server SERVER         Mail server [default: `+defaultserver+`].
        -P|--port PORT             Port, like 25 or 465 [default: `+defaultport+`]. ^5
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
       Commandline options take precedence over CONFIGFILE options.
    4. All given in the CONFIGFILE and on the commandline will be used.
    5. StartTLS is the default, except when PORT is 465, then SSL/TLS is used.
`, self, version)
}

func exitmsg(msg string) {
	usage()
	fmt.Printf("\nERROR %v\n", msg)
	os.Exit(1)
}

func readCsv(csvfile string) {
	file, err := os.Open(csvfile)
	if err != nil {
		exitmsg("Cannot open CSV file")
  }

  defer file.Close()
	reader := csv.NewReader(file)
	var headerread bool
	for {
		fields, err := reader.Read()
		if err == io.EOF { // End of records
			break
		}
		if err != nil {
			exitmsg("Error reading CSV file")
		}

		if headerread { // Record
			record := make(Record, len(csvheader))
			for i, key := range csvheader {
				record[key] = fields[i]
			}
			records = append(records, record)
		} else { // Header
			csvheader = fields
			headerread = true
		}
	}
}

func main() {
	var cfg Config
	var err error
	nArgs := len(os.Args)
	if nArgs == 1 {
		usage()
		return
	}
	// Use name of binary/link
	self = os.Args[0]
	i := strings.IndexByte(self, '/')
	for i >= 0 {
		self = self[i+1:]
		i = strings.IndexByte(self, '/')
	}

	// Parse commandline
	i = 1
	var from, to, subject, user, password, server, port, cc, bcc, reply, read, message, mfile, nmessage, nfile, cfile, csvfile string
	var ssltls bool
	var attachments []string
	for i < nArgs {
		switch os.Args[i] {
		case "-o", "--options":
			if cfile != "" {
				exitmsg("Can't use -o/--options twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -o/--options must have an argument")
			}

			cfile = os.Args[i+1]
			_, err := os.Stat(cfile)
			if err == nil { // File present
				cfgdata, err := os.ReadFile(cfile)
				if err != nil {
					exitmsg("Cannot read config file '" + cfile + "'")
				}

				err = yaml.UnmarshalStrict(cfgdata, &cfg)
				if err != nil {
					exitmsg("Error in config file '" + cfile + "':\n"+err.Error())
				}

			} else { // Not found
				exitmsg("Config file '" + cfile + "' not found")
			}

			i = i + 1
		case "-f", "--from":
			if from != "" {
				exitmsg("Can't use -f/--from twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -f/--from must have an argument")
			}

			from = os.Args[i+1]
			i = i + 1
		case "-t", "--to":
			if to != "" {
				exitmsg("Can't use -t/--to twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -t/--to must have an argument")
			}

			to = os.Args[i+1]
			i = i + 1
		case "-s", "--subject":
			if subject != "" {
				exitmsg("Can't use -s/--subject twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -s/--subject must have an argument")
			}

			subject = os.Args[i+1]
			i = i + 1
		case "-u", "--user":
			if user != "" {
				exitmsg("Can't use -u/--user twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -u/--user must have an argument")
			}

			user = os.Args[i+1]
			i = i + 1
		case "-p", "--password":
			if password != "" {
				exitmsg("Can't use -p/--password twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -p/--password must have an argument")
			}

			password = os.Args[i+1]
			i = i + 1
		case "-S", "--server":
			if server != "" {
				exitmsg("Can't use -S/--server twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -S/--server must have an argument")
			}

			server = os.Args[i+1]
			i = i + 1
		case "-P", "--port":
			if port != "" {
				exitmsg("Can't use -P/--port twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -P/--port must have an argument")
			}

			port = os.Args[i+1]
			i = i + 1
		case "-T", "--tls":
			ssltls = true
		case "-c", "--cc":
			if cc != "" {
				exitmsg("Can't use -c/--cc twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -c/--cc must have an argument")
			}

			cc = os.Args[i+1]
			i = i + 1
		case "-b", "--bcc":
			if bcc != "" {
				exitmsg("Can't use -b/--bcc twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -b/--bcc must have an argument")
			}

			bcc = os.Args[i+1]
			i = i + 1
		case "-r", "--reply":
			if reply != "" {
				exitmsg("Can't use -r/--reply twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -r/--reply must have an argument")
			}

			reply = os.Args[i+1]
			i = i + 1
		case "-R", "--read":
			if read != "" {
				exitmsg("Can't use -R/--read twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -R/--read must have an argument")
			}

			read = os.Args[i+1]
			i = i + 1
		case "-m", "--message":
			if message != "" {
				exitmsg("Can't use -m/--message twice")
			}

			if mfile != "" {
				exitmsg("Can't use both -m/--message and -M/--mfile flags")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -m/--message must have an argument")
			}

			message = os.Args[i+1]
			i = i + 1
		case "-M", "--mfile":
			if mfile != "" {
				exitmsg("Can't use -F/--file twice")
			}

			if message != "" {
				exitmsg("Can't use both -m/--message and -M/--mfile flags")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -M/--mfile must have an argument")
			}

			mfile = os.Args[i+1]
			i = i + 1
		case "-n", "--nmessage":
			if nmessage != "" {
				exitmsg("Can't use -n/--nmessage twice")
			}

			if nfile != "" {
				exitmsg("Can't use both -n/--nmessage and -N/--nfile flags")
			}
			if i+2 > len(os.Args) {
				exitmsg("Flag -n/--nmessage must have an argument")
			}

			nmessage = os.Args[i+1]
			i = i + 1
		case "-N", "--nfile":
			if nfile != "" {
				exitmsg("Can't use -N/--nfile twice")
			}

			if message != "" {
				exitmsg("Can't use both -n/--nmessage and -N/--nfile flags")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -N/--nfile must have an argument")
			}

			nfile = os.Args[i+1]
			i = i + 1
		case "-C", "--csv":
			if csvfile != "" {
				exitmsg("Can't use -C/--csv twice")
			}

			if i+2 > len(os.Args) {
				exitmsg("Flag -C/--csv must have an argument")
			}

			csvfile = os.Args[i+1]
			i = i + 1
		case "-a", "--attachment":
			if i+2 > len(os.Args) {
				exitmsg("Flag -a/--attachment must have an argument")
			}

			attachment := os.Args[i+1]
			if _, err := os.Stat(attachment); err == nil {
				attachments = append(attachments, attachment)
			} else {
				exitmsg("Attachment '"+attachment+"' not found")
			}

			i = i + 1
		case "-h", "--help":
			usage()
			return
		default:
			exitmsg("unknown commandline option: " + os.Args[i])
		}

		i += 1
	}
	if csvfile == "" && cfg.CSVfile != "" {
		csvfile = cfg.CSVfile
	}
	if csvfile != "" {
		readCsv(csvfile)
	}
	cfgtls := strings.ToLower(cfg.TLS)
	if cfgtls != "" && cfgtls != "0" && cfgtls != "no" && cfgtls != "false" && cfgtls != "off" {
		ssltls = true
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
	if read == "" {
		read = cfg.Read
	}
	if to == "" {
		to = cfg.To
	}
	if to == "" {
		exitmsg("Essential option 'to' missing")
	}

	if subject == "" {
		subject = cfg.Subject
	}
	if subject == "" {
		exitmsg("Essential option 'subject' missing")
	}

	if user == "" {
		user = cfg.User
	}
	if user == "" {
		exitmsg("Essential option 'user' missing")
	}

	if password == "" {
		password = cfg.Password
	}
	if password == "" {
		exitmsg("Essential option 'password' missing")
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
	if message == "" && mfile == "" { // Rely on configfile for plaintext body
		if cfg.Message != "" && cfg.Mfile != "" { // Both set
			exitmsg("Can't have both 'message' and 'mfile' options set in .mailer")
		}

		if cfg.Mfile != "" { // Mfile set, use it
			mfile = cfg.Mfile
		} else { // Message either set or empty
			message = cfg.Message
		}
	}
	if nmessage == "" && nfile == "" { // Rely on configfile for html body
		if cfg.Nmessage != "" && cfg.Nfile != "" { // Both set
			exitmsg("Can't have both 'nmessage' and 'nfile' options set in .mailer")
		}

		if cfg.Nfile != "" { // Nfile set, use it
			nfile = cfg.Nfile
		} else { // Message either set or empty
			nmessage = cfg.Nmessage
		}
	}
	if message == "" && mfile == "" && nmessage == "" && nfile == "" {
		exitmsg("Content missing, none of 'message'/'mfile'/'nmessage'/'nfile' given")
	}

	if cfg.Attachment != "" {
		if _, err := os.Stat(cfg.Attachment); err != nil {
			exitmsg("Attachment '"+cfg.Attachment+"' from .mailer not found")
		}

		attachments = append(attachments, cfg.Attachment)
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
			exitmsg("No 'from' email nor 'user' email")
		}
	}

	// Populate email
	if message == "" && mfile != "" {
		f, err := os.ReadFile(mfile)
		if err != nil {
			exitmsg("Mfile not found: '" + mfile + "'")
		}

		message = string(f)
	}
	if nmessage == "" && nfile != "" {
		f, err := os.ReadFile(nfile)
		if err != nil {
			exitmsg("Nfile not found: '" + nfile + "'")
		}

		nmessage = string(f)
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
	if read != "" {
		mail.ReadReceipt = strings.Split(read, ",")
	}
	mail.Subject = subject
	mail.Text = []byte(message)
	mail.HTML = []byte(nmessage)
	for _, attachment := range attachments {
		_, err := mail.AttachFile(attachment)
		if err != nil {
			exitmsg("Could not attach file '"+attachment+"'")
		}
	}

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
