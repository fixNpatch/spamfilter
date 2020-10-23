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
	Body           Body
}

type Address struct {
	From string
	To   string
}

type NominalAddress struct {
	From string
	TO   string
}

type Headers struct {
	Date      string
	MessageId string
	Subject   string

	MIME string

	AuthenticationResults string
	ReceivedSPF           string

	XHeaders
}

type XHeaders struct {
	XConfirmReadingTo string
	XErrorsTo         string
	XMailer           string
	XSender           string
	XPriority         string
	XUIDL             string
}

type Body struct {
	string
}


func returnValue(input string) (string, error) {
	// основанно на следующем знании: https://dmorgan.info/posts/encoded-word-syntax/

	// python version
	// encoded_word_regex = r'=\?{1}(.+)\?{1}([B|Q])\?{1}(.+)\?{1}='
	// charset, encoding, encoded_text = re.match(encoded_word_regex, encoded_words).groups()


	// регулярное выражение для определения кодировки в поле письма
	regular := regexp.MustCompile(`\?{1}(.+)\?{1}([B|Q])\?{1}(.+)\?{1}=`)
	splitted := regular.FindStringSubmatch(input) // выделение частей
	charst, encodng, encodtxt := splitted[1], splitted[2], splitted[3]

	// проверка (можно удалить вывод)
	fmt.Println(charst, encodng, encodtxt)

	// декодирование по выбранной кодировке (пока сделано только для base64)
	decoded, err := b64.StdEncoding.DecodeString("0KDQsNGB0YLQvtGA0LbQtdC90LjQtSDQv9C+0LTQv9C40YE=")
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


	// "Очистка данных
	for _, line := range rawLineArray {

		// если lineArray[index] содержит заголовок (проверяем при помощи regexp)
		// то добавляем эту строку к последнему заголовку

		match, err := regexp.MatchString("^[A-Z].[a-zA-Z-0-9]+:", line)
		if err != nil {
			fmt.Println("Error::Something wrong with your Regexp")
			return fmt.Errorf("Error:Cant parse file::Bad input")
		}


		if !match {
			modifiedLineArray[len(modifiedLineArray) - 1] =
				strings.TrimRight(string(modifiedLineArray[len(modifiedLineArray) - 1]) +
					strings.ReplaceAll(line, "\t", ""), "\r\n")
		} else {
			matched := strings.TrimRight(line, "\r\n")
			modifiedLineArray = append(modifiedLineArray, matched)
		}

	}


	// распаковка данных по структуре Mail (сопоставление заголовков с данными)
	for _, line := range modifiedLineArray {

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

	s, _ := returnValue("=?UTF-8?B?0KDQsNGB0YLQvtGA0LbQtdC90LjQtSDQv9C+0LTQv9C40YE=?=")
	fmt.Println(s)

	return nil
}

func (m *Mail_) Checker() {

	fmt.Println(m.NominalAddress.TO)
	fmt.Println(m.NominalAddress.From)
	fmt.Println(m.Address.From)

	//reasonsNumber=0

	//if len(m.Headers.Subject) > 120:
	//reasons

	return
}
