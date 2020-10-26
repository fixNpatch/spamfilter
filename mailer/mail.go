package mailer

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const Inbox = "mail/inbox"

type Mail_ struct {
	Address        Address
	NominalAddress NominalAddress
	Headers        Headers
	BodyParts      []string
}

type Address struct {
	From string
	To   string
}

type NominalAddress struct {
	From string
	To   string
}

type Headers struct {
	Date      string
	MessageId string
	Subject   string

	MIME string

	DKIMSignature []string

	AuthenticationResults string
	ReceivedSPF           string

	ReplyTo   string
	InReplyTo string

	XHeaders
}

// X-Заголовки не являются стандартными. В письме они могут быть указаны как угодно, могут быть новыми.
// Поэтому выделено наиболее часто употребимые
type XHeaders struct {
	XConfirmReadingTo string // запрашивает автоматическое подтверждение того, что письмо было получено или прочитано. Предполагается соответствующая реакция почтовой программы, но обычно он игнорируется.
	XErrorsTo         string
	XMailer           string
	XSender           string
	XPriority         string
	XSpam             string
	XUIDL             string
}

func returnValue(input string) (string, error) {
	// основанно на следующем знании: https://dmorgan.info/posts/encoded-word-syntax/

	// регулярное выражение для определения кодировки в поле письма
	regular := regexp.MustCompile(`\?{1}(.+)\?{1}([B|Q])\?{1}(.+)\?{1}=`)
	splitted := regular.FindStringSubmatch(input) // выделение частей
	charst, encodng, encodtxt := splitted[1], splitted[2], splitted[3]

	_ = encodng + charst // заглушка

	// декодирование по выбранной кодировке (пока сделано только для base64)
	decoded, err := b64.StdEncoding.DecodeString(encodtxt)
	if err != nil {
		fmt.Println("Error::Bad converting::Base64 => UTF-8")
		return "", err
	}
	return string(decoded), nil
}

func readfile(filename string) string {
	data, err := ioutil.ReadFile(Inbox + "/" + filename)
	if err != nil { // если что-то пошло не так, то сохраняем ошибку и возвращаем пустую строку
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func (m *Mail_) Parser(info os.FileInfo) error {
	var modifiedLineArray []string // массив строк, куда будем складывать полностью форматированные строки

	// входные данные
	readed := readfile(info.Name())
	if readed == "" { // если данный файл пустой, то выдаем ошибку
		fmt.Println("Error::Cant read file")
		return fmt.Errorf("Error:Cant read file::Empty input")
	}

	// regexp для заголовков ^[a-zA-Z].[a-zA-Z-0-9]+:

	rawLineArray := strings.Split(strings.TrimSuffix(strings.TrimSuffix(readed, "\n"), "\r"), "\n")
	// разбиваем считанную из файла информацию на строки
	// получаем "грязные" строки, то есть некоторые заголовки являются многострочными.
	// Надо сделать 1 заголовок = 1 строка

	// "Очистка" данных
	for _, line := range rawLineArray {

		// если lineArray[index] содержит заголовок (проверяем при помощи regexp)
		// то добавляем эту строку к последнему заголовку

		match, err := regexp.MatchString("^[A-Z].[a-zA-Z-0-9]+:", line)
		if err != nil {
			fmt.Println("Error::Something wrong with your Regexp")
			return fmt.Errorf("Error:Cant parse file::Bad input")
		}

		if !match {
			modifiedLineArray[len(modifiedLineArray)-1] =
				strings.TrimRight(string(modifiedLineArray[len(modifiedLineArray)-1])+
					strings.ReplaceAll(line, "\t", ""), "\r\n")
		} else {
			matched := strings.TrimRight(line, "\r\n")
			modifiedLineArray = append(modifiedLineArray, matched)
		}

	}

	var clearLine string

	// распаковка данных по структуре Mail (сопоставление заголовков с данными)
	for _, line := range modifiedLineArray {

		// проверяем есть декодированные части в очередной рассматриваемой строке
		match, err := regexp.MatchString(`\?{1}(.+)\?{1}([B|Q])\?{1}(.+)\?{1}=`, line)
		if err != nil {
			fmt.Println("Error::Something wrong with your Regexp")
			return fmt.Errorf("Error:Cant parse file::Bad input")
		}
		if match { // если есть закодированные части, то необходимо раскодировать каждую из них
			sublines := strings.Split(line, " ") // берем разделение по пробелам
			for _, subline := range sublines {   // рассматриваем каждую из частей
				// если взята закодированная часть, то декодируем её и результат добавляем в результирующую строку
				if match, err := regexp.MatchString(`\?{1}(.+)\?{1}([B|Q])\?{1}(.+)\?{1}=`, subline); err == nil && match {
					tmp, err := returnValue(subline)
					if err != nil {
						fmt.Println(err)
						return err
					}
					clearLine += tmp + " "
				} else { // если часть не закодированна, то просто добавляем её в результат
					clearLine += subline + " "
				}
			}
			line = clearLine // заменяем оригинальную строку результирующей.
		}

		if len(line) >= 15 && line[:12] == "Delivered-To" {
			m.NominalAddress.To = line[14:]
		}

		if len(line) >= 5 && line[:2] == "To" {
			m.Address.To = line[4:]
		}

		if len(line) >= 14 && line[:11] == "Return-Path" {
			m.NominalAddress.From = line[13:]
		}

		if len(line) >= 7 && line[:4] == "From" {
			m.Address.From = line[6:]
		}

		if len(line) >= 25 && line[:22] == "Authentication-Results" {
			m.Headers.AuthenticationResults = line[24:]
		}

		if len(line) >= 15 && line[:12] == "Received-SPF" {
			m.Headers.ReceivedSPF = line[14:]
		}

		if len(line) >= 11 && line[:8] == "Received" {
			m.Address.From = line[10:]
		}

		if len(line) >= 13 && line[:10] == "Message-ID" {
			m.Headers.MessageId = line[12:]
		}

		if len(line) >= 10 && line[:7] == "Subject" {
			m.Headers.Subject = line[9:]
		}

		if len(line) >= 17 && line[:14] == "DKIM-Signature" {
			m.Headers.DKIMSignature = append(m.Headers.DKIMSignature, line[16:])
		}

		if len(line) >= 7 && line[:4] == "Date" {
			m.Headers.Date = line[6:]
		}

		if len(line) >= 11 && line[:8] == "Reply-To" {
			m.Headers.ReplyTo = line[10:]
		}

		if len(line) >= 14 && line[:11] == "In-Reply-To" {
			m.Headers.ReplyTo = line[13:]
		}

	}

	return nil
}
