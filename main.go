package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"spamfilter/mailer"
	"strings"
)


func main()  {
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

		m := new(mailer.Mail_)

		// Для каждого отдельного файла из списка этих файлов
		for _, file := range files {
			fmt.Println(file.Name()) // Выводим все имена файлов
			filter(file, m)
		}

		fmt.Println("Все письма отсортированы")

	} else {
		fmt.Println("Нет новых писем")
	}
}

func check_(m *mailer.Mail_){

	// Контрольная сумма >= 4 очков => спам.
	// Признаки в заголовках оцениваются в 1,40
	// Признаки в теле письма оцениваются в 2
	// Следовательно необходимо 3/0 или 2/1 или 0/2


	checkSum := 0.0

	// 1. Сравнить слова в теме и тексте письма: должны быть совпадения

	// 2. Title(Subject) документа длинней 12 слов или 120 символов.
	splice := strings.Split(m.Headers.Subject, " ")
	if len(splice) > 12 || len(m.Headers.Subject) > 120 {
		fmt.Println("Проблема с заголовком")
		checkSum += 1.40
	}

	// 3. Сообщение разослано большому количеству пользователей

	// 4. Отсутствие обратного адреса




	// ----------------------- По телу --------------------------------

	// 1. Meta name="description" длинней 40 слов или 250 символов. (ЕСЛИ HTML ФОРМАТ)

	// 2. Meta name="keywords" длинней 40 слов или 250 символов. (ЕСЛИ HTML ФОРМАТ)

	// 3. Проверка ссылок (VirusTotal Api) (сразу +4 влепливаем)

	// 4. Количество тегов на количество текста (Предполагаю не более 1 тега на 10 слов)

	// 5. Только картинка


}


func filter(info os.FileInfo, m *mailer.Mail_){
	fmt.Println("Найдено новое письмо:", info.Name())

	err := m.Parser(info) // Распределяем данные по структуре письма
	if err != nil {
		fmt.Println(err)
		return
	}


	check_(m) // Производим проверку

	return
}