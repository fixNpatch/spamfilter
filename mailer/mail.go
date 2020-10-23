package mailer

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const Inbox = "mail/inbox"

type Mail_ struct {
	Address Address
	NominalAddress NominalAddress
	Headers Headers
	Body Body
}

type Address struct {
	From string
	To string
}

type NominalAddress struct {
	From string
	TO string
}

type Headers struct {
	Date string
	MessageId string
	Subject string

	MIME string

	AuthenticationResults string
	ReceivedSPF string

	XHeaders
}

type XHeaders struct {
	XConfirmReadingTo string
	XErrorsTo string
	XMailer string
	XSender string
	XPriority string
	XUIDL string
}

type Body struct {
	string
}

func readfile(filename string) string {
	data, err := ioutil.ReadFile(Inbox + "/" + filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func (m *Mail_) Parser(info os.FileInfo)(error)  {
	readed := readfile(info.Name())
	if readed == "" {
		fmt.Println("Error::Cant read file")
		return fmt.Errorf("Error:Cant read file::Empty input")
	}


	for _, line := range strings.Split(strings.TrimSuffix(readed, "\n"), "\n") {

		if len(line) >= 15 && line[:12] == "Delivered-To" {
			m.NominalAddress.TO = line[14:]
		}

		if len(line) >= 14 && line[:11] == "Return-path" {
			m.NominalAddress.From = line[13:]
		}

		if len(line) >= 25 && line[:22] == "Authentication-Results" {
			m.Address.From = line[24:]
		}


		if len(line) >= 7 && line[:4] == "From" {
			m.Address.From = line[6:]
		}

		if len(line) >= 5 && line[:2] == "From" {
			m.Address.From = line[4:]
		}

		if len(line) >= 12 && line[:10] == "Message-ID" {
			m.Headers.MessageId = line[11:]
		}

		if len(line) >=9 && line[:7] == "Subject" {
			m.Headers.Subject = line[8:]
		}


	}


	//fmt.Println(readed)

	re := regexp.MustCompile(`To`)
	fmt.Printf("%q\n", re.FindStringSubmatch(readed))
	//m.Body.string = "Body"

	return nil
}

func (m *Mail_) Checker()  {

	fmt.Println(m.NominalAddress.TO)
	fmt.Println(m.NominalAddress.From)
	fmt.Println(m.Address.From)

	//reasonsNumber=0

	//if len(m.Headers.Subject) > 120:
		//reasons


	return
}