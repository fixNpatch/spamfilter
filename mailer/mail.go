package mailer

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"spamfilter/configurator"
	"strings"
)

type Mail_ struct {
	Address        Address
	NominalAddress NominalAddress
	Headers        Headers
	BodyParts      []string
}

// Указательный адрес. Можем писать туда что угодно. Для удобства пользователей
type Address struct {
	From string
	To   string
}

// Конвертный адрес. То, что нужно, чтобы доставить письмо. Технически важен
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

func readfile(filename string, cfg *configurator.Config) string {
	data, err := ioutil.ReadFile(cfg.InboxPath + "\\" + filename)
	if err != nil { // если что-то пошло не так, то сохраняем ошибку и возвращаем пустую строку
		fmt.Println("Error::Mailer::readfile::ReadFile::", err)
		return ""
	}
	return string(data)
}

func (m *Mail_) Parser(info os.FileInfo, cfg *configurator.Config) error {
	var modifiedLineArray []string // массив строк, куда будем складывать полностью форматированные строки

	// входные данные
	readed := readfile(info.Name(), cfg)
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

		match, err := regexp.MatchString("^[A-Z][a-zA-Z-0-9]+:", line)
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

		clearLine = ""
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
						fmt.Println("Error::Mailer::?::", err)
						return err
					}
					clearLine += strings.TrimSuffix(tmp, " ") + " "
				} else { // если часть не закодированна, то просто добавляем её в результат
					clearLine += subline + " "
				}
			}
		}

		if clearLine == "" {
			clearLine = line
		}

		if len(clearLine) >= 15 && strings.ToLower(clearLine[:12]) == "delivered-to" {
			m.NominalAddress.To = clearLine[14:]
		}

		if len(clearLine) >= 5 && strings.ToLower(clearLine[:2]) == "to" {
			m.Address.To = clearLine[4:]
		}

		if len(clearLine) >= 14 && strings.ToLower(clearLine[:11]) == "return-path" {
			m.NominalAddress.From = clearLine[13:]
		}

		if len(clearLine) >= 7 && strings.ToLower(clearLine[:4]) == "from" {
			m.Address.From = clearLine[6:]
		}

		if len(clearLine) >= 25 && strings.ToLower(clearLine[:22]) == "authentication-results" {
			m.Headers.AuthenticationResults = clearLine[24:]
		}

		if len(clearLine) >= 15 && strings.ToLower(clearLine[:12]) == "received-spf" {
			m.Headers.ReceivedSPF = clearLine[14:]
		}

		if len(clearLine) >= 11 && strings.ToLower(clearLine[:8]) == "received" {
			m.Address.From = clearLine[10:]
		}

		if len(clearLine) >= 13 && strings.ToLower(clearLine[:10]) == "message-id" {
			m.Headers.MessageId = clearLine[12:]
		}

		if len(clearLine) >= 10 && strings.ToLower(clearLine[:7]) == "subject" {
			m.Headers.Subject = clearLine[9:]
		}

		if len(clearLine) >= 17 && strings.ToLower(clearLine[:14]) == "dkim-signature" {
			m.Headers.DKIMSignature = append(m.Headers.DKIMSignature, clearLine[16:])
		}

		if len(clearLine) >= 7 && strings.ToLower(clearLine[:4]) == "date" {
			m.Headers.Date = clearLine[6:]
		}

		if len(clearLine) >= 11 && strings.ToLower(clearLine[:8]) == "reply-to" {
			m.Headers.ReplyTo = clearLine[10:]
		}

		if len(clearLine) >= 14 && strings.ToLower(clearLine[:11]) == "in-reply-to" {
			m.Headers.InReplyTo = clearLine[13:]
		}

	}

	return nil
}
