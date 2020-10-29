package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"spamfilter/mailer"
	"strings"
)

var MyEmail string
var CurrentDir string
var ListOfMails []string

const Inbox = "\\mail\\inbox"
const SpamBox = "\\mail\\bad"
const ClearBox = "\\mail\\good"

func main() {
	var err error
	fmt.Println(">> Фильтр запущен <<")

	CurrentDir, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	MyEmail = "titov.ant.workmail@gmail.com"
	fmt.Println("Введите свой email (будет использован для проверки)")
	_, err = fmt.Scanf("%s\n", &MyEmail)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Смотрим наличие файлов во входящих
	files, err := ioutil.ReadDir(mailer.Inbox)
	if err != nil {
		log.Fatal(err)
	}
	// если есть файлы в рассматриваемой выше директории, то...
	if len(files) > 0 {

		/*TODO do not delete following*/

		//a := mailer.Address{From:"", To:""}
		//x := mailer.XHeaders{"", "","","","",""}
		//h := mailer.Headers{"","","","","","","", x}
		//b := mailer.Body{""}
		//m := mailer.Mail_{Address:a, Headers:h} // создаем Менеджер писем

		// Для каждого отдельного файла из списка этих файлов
		for _, file := range files {

			m := new(mailer.Mail_)

			filter(file, m)     // проводим процедуру фильтрации
			fmt.Println("\n\n") // делаем отступ в консоли
		}

		fmt.Println("Все письма отсортированы")

	} else {
		fmt.Println("Нет новых писем")
	}
}

// данный функционал будет реализован на дипломной работе
func check_headers(m *mailer.Mail_) float64 {
	var splice []string

	checkSum := 0.0

	// 1. Сравнить слова в теме и тексте письма: должны быть совпадения

	// 2. Title(Subject) документа длинней 12 слов или 120 символов.
	splice = strings.Split(strings.TrimPrefix(m.Headers.Subject, " "), " ")
	if len(splice) > 12 || len(splice) <= 1 || len(m.Headers.Subject) > 120 {
		fmt.Println("Внимание: Проблема с заголовком")
		fmt.Println(splice)
		checkSum += 1.40
	}

	// 3. Сообщение разослано большому количеству пользователей

	splice = strings.Split(strings.TrimPrefix(m.Address.To, " "), " ")
	if len(splice) > 12 || len(splice) <= 1 || len(m.Headers.Subject) > 120 {
		fmt.Println("Внимание: Проблема с получателями")
		fmt.Println(splice)
		checkSum += 1.40
	}

	fmt.Println(m.Address.From)
	fmt.Println(m.NominalAddress.From)

	// 4. Отсутствие обратного адреса
	// 5. Проверка DKIM. Не реализовано. Оставил на дипломную работу
	// что почитать https://securelist.ru/texnologiya-dkim-na-strazhe-vashej-pochty/25010/

	return checkSum
}

// данный функционал будет реализован на дипломной работе
func check_body(m *mailer.Mail_) float64 {
	checkSum := 0.0

	// ----------------------- По телу --------------------------------

	// 1. Meta name="description" длинней 40 слов или 250 символов. (ЕСЛИ HTML ФОРМАТ)
	// 2. Meta name="keywords" длинней 40 слов или 250 символов. (ЕСЛИ HTML ФОРМАТ)
	// 3. Проверка ссылок (VirusTotal Api) (сразу +4 влепливаем)
	// 4. Количество тегов на количество текста (Предполагаю не более 1 тега на 10 слов)
	// 5. Только картинка

	return checkSum
}

func check_(m *mailer.Mail_) bool {

	var splice []string

	// Контрольная сумма >= 4 очков => спам.
	// Признаки в заголовках оцениваются в 1.40
	// Признаки в теле письма оцениваются в 2.00
	// Следовательно необходимо 3/0 или 2/1 или 0/2

	fmt.Println("---------- Проверка данного письма  ---------------")

	checkSum := 0.0

	// следующие проверки будут реализованы в дипломной работе
	// checkSum += check_headers(m)
	// checkSum += check_body(m)

	// ------------------ (что нужно сдать) --------------------

	// 0. Заголовок письма
	splice = strings.Split(strings.TrimPrefix(m.Headers.Subject, " "), " ")
	if len(splice) > 12 || len(splice) <= 1 || len(m.Headers.Subject) > 120 {
		fmt.Println("Внимание: Проблема с заголовком")

		for i, item := range splice {
			fmt.Println("\t", i, item)
		}

		checkSum += 1.40
	}

	// 1. Обратный адрес Если пустой то +rating
	if len(m.NominalAddress.From) < 1 {
		fmt.Println("Внимание: Проблема с обратным адресом")
		fmt.Println("\tНе указан или отсутствует заголовок Return-Path")
		checkSum += 1.40
	}

	// 2. Получатели: если пустой или их много, то +rating
	splice = strings.Split(strings.TrimPrefix(m.NominalAddress.To, " "), " ")
	if len(splice) > 5 || len(splice) < 1 {
		fmt.Println("Внимание: Проблема с количеством получателей")
		fmt.Println("\tКоличество получателей:", len(splice))
		for i, item := range splice {
			fmt.Println("\t", i, item)
		}
		checkSum += 1.40
	}

	// 3. Есть ли наш адрес в получателях. Если нет, то +rating
	splice = strings.Split(strings.TrimPrefix(m.Address.To, " "), " ")
	found := false
	for _, addr := range splice {
		if len(addr) > 3 && addr[1:len(addr)-1] == MyEmail {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Внимание: Проблема с неверно указанным получателем")

		for i, item := range splice {
			if item != "" {
				fmt.Println("\t", i, item, "!=", MyEmail)
			}
		}

		checkSum += 1.40
	}

	// 4. Ищем DKIM сигнатуру. Если её нет, то +rating
	if len(m.Headers.DKIMSignature) < 1 {
		fmt.Println("Внимание: Проблема с DKIM")
		if len(m.Headers.DKIMSignature) == 0 {
			fmt.Println("\tПустой заголовок DKIM")
		}
		checkSum += 2.00
	}

	// 5. Reply-To. Если пустой то +rating
	if len(m.Headers.ReplyTo) < 1 {
		fmt.Println("Внимание: Проблема с Reply-To")
		if m.Headers.ReplyTo == "" {
			fmt.Println("\tОтсутствует заголовок Reply-To")
		}
		checkSum += 1.40
	}

	// 6. In-Reply-To. Если пустой то +rating
	if len(m.Headers.InReplyTo) < 1 {
		fmt.Println("Внимание: Проблема с In-Reply-To")
		if m.Headers.InReplyTo == "" {
			fmt.Println("\tОтсутствует заголовок In-Reply-To")
		}
		checkSum += 1.00
	}

	// 7. Message-ID. Если пустой то +rating
	if len(m.Headers.MessageId) < 1 {
		fmt.Println("Внимание: Проблема с Message-Id")
		if m.Headers.MessageId == "" {
			fmt.Println("\tОтсутствует заголовок Message-Id")
		}
		checkSum += 1.40
	}

	// 8. Received-SPF. Если пустой то +rating
	if len(m.Headers.ReceivedSPF) < 1 {
		fmt.Println("Внимание: Проблема с Received-SPF")
		if m.Headers.ReceivedSPF == "" {
			fmt.Println("\tОтсутствует заголовок Received-SPF")
		}
		checkSum += 1.40
	}

	if checkSum < 2.0 {
		fmt.Println("---------------------------------------------------")
		fmt.Println("\t\tПисьмо точно не спам")
		fmt.Println("---------------------------------------------------")
		return false
	} else if checkSum < 4.0 {
		fmt.Println("---------------------------------------------------")
		fmt.Println("\t\tПисьмо скорее всего не спам")
		fmt.Println("---------------------------------------------------")
		return false
	} else {
		fmt.Println("---------------------------------------------------")
		fmt.Println("\t\tОбнаружен спам!")
		fmt.Println("---------------------------------------------------")
		return true
	}
}

func filter(info os.FileInfo, m *mailer.Mail_) {
	fmt.Println("Найдено новое письмо:", info.Name())

	err := m.Parser(info) // Распределяем данные по структуре письма
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(info.Name())
	action := check_(m) // Производим проверку

	if action {
		err = os.Rename(CurrentDir+Inbox+"\\"+info.Name(), CurrentDir+SpamBox+"\\"+info.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err = os.Rename(CurrentDir+Inbox+"\\"+info.Name(), CurrentDir+ClearBox+"\\"+info.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}
