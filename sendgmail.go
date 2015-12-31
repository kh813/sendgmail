// sendgmail.go
// last updated : 2015-12-31

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

const (
	//version = ""
	default_smtpauthuser = "" // Default(Your) Gmail account
	default_smtpauthpass = "" // Default(Your) Gmail password
)

func setMailfrom(from string, smtpauthuser string) string {
	if from != "" {
		//fmt.Println("setting from from CLI value")     // debug
		return from
	} else {
		//fmt.Println("setting from from smtpauthuser")  // debug
		return smtpauthuser
	}
}

func setMaildest(to string, cc string) []string {
	if cc == "" {
		maildest := []string{to}
		return maildest
	} else {
		maildest := []string{to, cc}
		return maildest
	}
}

// ref :  http://qiita.com/yamasaki-masahide/items/a9f8b43eeeaddbfb6b44
func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

func utf8Split(utf8string string, length int) []string {
	resultString := []string{}
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

func encodeSubject(subject string) string {
	if subject != "" {
		var buffer bytes.Buffer
		buffer.WriteString("Subject: ")
		for _, line := range utf8Split(subject, 13) {
			buffer.WriteString(" =?utf-8?B?")
			buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
			buffer.WriteString("?=\r\n")
		}
		return buffer.String()
	} else {
		// \r\n is required to avoid the next line set as subject,
		// when the password is empty
		return "Subject: \r\n"
	}
}

type Config struct {
	Gmailuser string
	Gmailpass string
}

func main() {
	// Get options CLI
	to := flag.String("t", "", "-t 'mail to'")
	cc := flag.String("c", "", "-c 'mail cc'")
	from := flag.String("f", "", "-f 'mail from'")
	subject := flag.String("s", "", "-s 'subject'")
	smtpauthuser := flag.String("g", "", "-g 'Gmail account email is sent from'")
	smtpauthpass := flag.String("p", "", "-p 'Gmail password'")
	//quiet := flag.Bool("q", false, "-q Quiet")
	verbose := flag.Bool("v", false, "-v Verbose")
	rawmode := flag.Bool("rm", false, "-rm Raw mode")
	//version := flag.Bool("v", false, "-v Show version")
	//debug := flag.Bool("d", false, "-d Debug mode")
	flag.Parse()

	// Default values
	smtphost := "smtp.gmail.com"
	smtpport := "587"
	smtpserver := smtphost + ":" + smtpport
	//maildest := []string{*to}
	maildest := setMaildest(*to, *cc)
	mailfrom := setMailfrom(*from, *smtpauthuser)

	// Load conf
	var conf Config
	var file []byte
	file, err := ioutil.ReadFile("sendgmail.conf")
	if err != nil {
		//fmt.Println("sendgmail.conf not found")
		if *verbose == true {
			fmt.Println("sendgmail.conf not found")
		}
		//panic(err)
	}
	if len(file) > 0 {
		jsonerr := json.Unmarshal(file, &conf)
		if jsonerr != nil {
			fmt.Println("sendgmail.conf found but cannot be read")
			/* if *quiet != true {
				fmt.Println("sendgmail.conf found but cannot be read")
			} */
			panic(jsonerr)
		}
	}

	// Check/set smtp auth-user/password
	if *smtpauthuser != "" && *smtpauthpass != "" {
		if *verbose == true {
			fmt.Println("gmail account from cli used")
		}
	} else if conf.Gmailuser != "" {
		*smtpauthuser = conf.Gmailuser
		*smtpauthpass = conf.Gmailpass
		if *verbose == true {
			fmt.Println("Gmail account from conf used")
		}
	} else if *smtpauthuser == "" && default_smtpauthuser != "" && default_smtpauthpass != "" {
		*smtpauthuser = default_smtpauthuser
		*smtpauthpass = default_smtpauthpass
		if *verbose == true {
			fmt.Println("Hardcoded Gmail account/password used")
		}
	}

	// Error checks
	if *to == "" {
		fmt.Println("Mail-to address is empty. Set it with -t option: -t 'account@example.com'")
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	if *smtpauthuser == "" || *smtpauthpass == "" {
		fmt.Println("Your Gmail user and/or password is empty..\nSet -g 'Gmail account' & -p 'smtp auth password'")
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	if *subject == "" {
		if *verbose == true {
			fmt.Println("Mail subject is empty, but sending it, anyway.\nSet -s 'subject' next time! ")
		}
	}

	// Set mail header
	mailheader := ""
	if *rawmode == false {
		mailheader = "" +
			"From: " + *smtpauthuser + "\r\n" +
			"To: " + *to + "\r\n" +
			"Cc: " + *cc + "\r\n" +
			//"Subject: " + *subject + "\r\n" +
			encodeSubject(*subject) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"Content-Transfer-Encoding: base64\r\n" +
			"\r\n"
	}

	// Get mail body or set default message
	mailbody := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		mailbody = mailbody + scanner.Text() + "\r\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	} else {
		//fmt.Println("Scanned successfullyl!")  // debug
	}

	// Set mail head + "encoded" mail body
	mailmsg := ""
	if *rawmode == false {
		mailmsg = mailheader + add76crlf(base64.StdEncoding.EncodeToString([]byte(mailbody)))
	} else {
		mailmsg = mailbody
	}
	// Version, debug, etc.
	/* if *debug == true {
		fmt.Println("Debug mode")
		fmt.Println("From :    " + *to)
		fmt.Println("CC   :    " + *cc)
		fmt.Println("Subject : " + *subject)
		fmt.Println("Auth user : " + *smtpauthuser)
		fmt.Println("Auth pass : " + *smtpauthpass)
		fmt.Println("Mail message : \n" + mailmsg)
		os.Exit(0)
	} else if *version == true {
		fmt.Println("Version : " + version)
		os.Exit(0)
	} */

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		*smtpauthuser,
		*smtpauthpass,
		smtphost,
	)

	// Connect to the server, authenticate, set the sender and recipient,
	smtperr := smtp.SendMail(
		smtpserver,
		auth,
		mailfrom,
		maildest,
		[]byte(mailmsg),
	)
	if smtperr != nil {
		log.Fatal(smtperr)
	} else {
		if *verbose == true {
			fmt.Println("Mail sent successfully")
		}
	}
}
