package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"spamfilter/mailer"
	"strings"
)

func main() {
	fmt.Println(">> Фильтр запущен <<")

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

			fmt.Println(file.Name()) // Выводим имя рассматриваемого файла
			filter(file, m)          // проводим процедуру фильтрации
			fmt.Println("\n\n")      // делаем отступ в консоли
		}

		fmt.Println("Все письма отсортированы")

	} else {
		fmt.Println("Нет новых писем")
	}
}

func check_(m *mailer.Mail_) {
	var splice []string

	// Контрольная сумма >= 4 очков => спам.
	// Признаки в заголовках оцениваются в 1.40
	// Признаки в теле письма оцениваются в 2.00
	// Следовательно необходимо 3/0 или 2/1 или 0/2

	fmt.Println("---------- MAILCHEKER for this mail ---------------")

	checkSum := 0.0

	// 1. Сравнить слова в теме и тексте письма: должны быть совпадения

	// 2. Title(Subject) документа длинней 12 слов или 120 символов.
	splice = strings.Split(strings.TrimPrefix(m.Headers.Subject, " "), " ")
	if len(splice) > 12 || len(splice) <= 1 || len(m.Headers.Subject) > 120 {
		fmt.Println("CAUTION: Проблема с заголовком")
		checkSum += 1.40
	}

	// 3. Сообщение разослано большому количеству пользователей

	splice = strings.Split(strings.TrimPrefix(m.Address.To, " "), " ")
	if len(splice) > 12 || len(splice) <= 1 || len(m.Headers.Subject) > 120 {
		fmt.Println("CAUTION: Проблема с получателями")
		checkSum += 1.40
	}

	fmt.Println(m.Address.From)
	fmt.Println(m.NominalAddress.From)

	// 4. Отсутствие обратного адреса

	// 5. Проверка DKIM. Не реализовано. Оставил на дипломную работу
	// что почитать https://securelist.ru/texnologiya-dkim-na-strazhe-vashej-pochty/25010/

	// ----------------------- По телу --------------------------------

	// 1. Meta name="description" длинней 40 слов или 250 символов. (ЕСЛИ HTML ФОРМАТ)

	// 2. Meta name="keywords" длинней 40 слов или 250 символов. (ЕСЛИ HTML ФОРМАТ)

	// 3. Проверка ссылок (VirusTotal Api) (сразу +4 влепливаем)

	// 4. Количество тегов на количество текста (Предполагаю не более 1 тега на 10 слов)

	// 5. Только картинка

	// ------------------ Изичи (что нужно сдать) --------------------

	// 1. Обратный адрес Если пустой то +rating
	// 2. Получатели: если пустой или их много, то +rating
	// 3. Есть ли наш адрес в получателях. Если нет, то +rating
	// 4. Ищем DKIM сигнатуру. Если её нет, то +rating
	// 5. Reply-To. Если пустой то +rating
	// 6. In-Reply-To. Если пустой то +rating
	// 7. Message-ID. Если пустой то +rating
	// 8. Received-SPF. Если пустой то +rating

	fmt.Println("---------------------------------------------------")
	fmt.Println("---------------------------------------------------")
}

func filter(info os.FileInfo, m *mailer.Mail_) {
	fmt.Println("Найдено новое письмо:", info.Name())

	err := m.Parser(info) // Распределяем данные по структуре письма
	if err != nil {
		fmt.Println(err)
		return
	}

	check_(m) // Производим проверку

	return
}
