package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"spamfilter/mailer"
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
	fmt.Println(m.Headers.Subject)
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