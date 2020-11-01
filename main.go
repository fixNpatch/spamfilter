package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"spamfilter/configurator"
	"spamfilter/mailer"
	"strings"
)

var cfg configurator.Config

var MyEmail string
var ListOfMails []string

var AbsoluteInbox, AbsoluteGood, AbsoluteSpam string

// Загрузка настроек
func configCheck(cfg *configurator.Config) (bool, error) {
	var err error
	if cfg.IsUsed {
		fmt.Println("Были найдены настроки")
		fmt.Println("Входящие:\t", cfg.InboxPath)
		fmt.Println("НеСпам:\t\t", cfg.FilteredPath)
		fmt.Println("Спам:\t\t", cfg.SpamPath)
		fmt.Println("Получатель:\t", cfg.TargetEmail)
		if len(cfg.ListOfMails) < 1 {
			fmt.Println("Список получателей не установлен")
		} else {
			fmt.Println("Список получателей:")
			for i, item := range cfg.ListOfMails {
				fmt.Println("\t", i, item)
			}
		}

		var changeCommand string
		for true {
			fmt.Println("Хотите изменить настройки? (y/N)")
			_, err = fmt.Scanf("%s\n", &changeCommand)
			if err != nil {
				fmt.Println(err)
				return false, err
			}
			if strings.ToLower(changeCommand) == "y" {
				return true, nil
			} else if strings.ToLower(changeCommand) == "n" {
				return false, nil
			} else {
				fmt.Println("Некорректный ввод. Пожалуйста введите символ Y(да) либо N(нет)")
				continue
			}
		}
	}
	return false,
		errors.New("ВНИМАНИЕ: Вы попали на участок кода, непредусмотренный алгоритмом. Свяжитесь с разработчиком")
}

// Сохранение настроек
func configSave(cfg *configurator.Config, master *configurator.Configurator) error {
	var err error
	var readyCheck string
	var ready bool
	for true {

		fmt.Println("Выбрано действие редактирования конфигурации =>\n" +
			"Если не требуется изменений, оставляйте пустую строку и нажимайте Enter.")

		fmt.Println("Введите свой email (будет использован для проверки).")
		if _, err = fmt.Scanf("%s\n", &MyEmail); err != nil &&
			err.Error() != "unexpected newline" {
			fmt.Println(err)
			return err
		}

		if MyEmail == "" {
			MyEmail = cfg.TargetEmail
		}
		fmt.Println("Получатель:", MyEmail)

		fmt.Println("Укажите абсолютный путь до директории входящих Email")
		if _, err = fmt.Scanf("%s\n", &AbsoluteInbox); err != nil &&
			err.Error() != "unexpected newline" {
			fmt.Println(err)
			return err
		}
		if AbsoluteInbox == "" {
			AbsoluteInbox = cfg.InboxPath
		}
		fmt.Println("Входящие:", AbsoluteInbox)

		fmt.Println("Укажите абсолютный путь до директории отфильтрованных сообщений")
		if _, err = fmt.Scanf("%s\n", &AbsoluteGood); err != nil &&
			err.Error() != "unexpected newline" {
			fmt.Println(err)
			return err
		}
		if AbsoluteGood == "" {
			AbsoluteGood = cfg.FilteredPath
		}
		fmt.Println("НеСпам:", AbsoluteGood)

		fmt.Println("Укажите абсолютный путь до директории спам-писем")
		if _, err = fmt.Scanf("%s\n", &AbsoluteSpam); err != nil &&
			err.Error() != "unexpected newline" {
			fmt.Println(err)
			return err
		}
		if AbsoluteSpam == "" {
			AbsoluteSpam = cfg.SpamPath
		}
		fmt.Println("Спам:", AbsoluteSpam)

		ready = false
		fmt.Println("Всё верно? (Y/n)")
		for true {
			if _, err = fmt.Scanf("%s\n", &readyCheck); err != nil && err.Error() != "unexpected newline" {
				fmt.Println(err)
				return err
			}
			if strings.ToLower(readyCheck) == "y" || readyCheck == "" {
				ready = true
				break
			} else if strings.ToLower(readyCheck) == "n" {
				break
			} else {
				fmt.Println("Некорректный ввод. Пожалуйста введите символ Y(да) либо N(нет)")
			}
		}
		if !ready {
			continue
		}

		newCfg := configurator.Config{
			IsUsed:       true,
			InboxPath:    AbsoluteInbox,
			FilteredPath: AbsoluteGood,
			SpamPath:     AbsoluteSpam,
			TargetEmail:  MyEmail,
			ListOfMails:  ListOfMails,
		}

		err = master.SetConfig(&newCfg)
		if err != nil {
			fmt.Println(err)
			return err
		}
		break
	}
	return nil
}

func main() {
	var err error
	fmt.Println(">> Фильтр запущен <<")

	configMaster := new(configurator.Configurator)

	// получаем настройки
	cfg, err = configMaster.GetConfig()
	if err != nil {
		fmt.Println("Error::Main::GetConfig::", err)
		return
	}

	// если настроки не были установлены
	if !cfg.IsUsed {
		err = configSave(&cfg, configMaster) // то производим первоначальную установку
	} else { // в противном случае
		changeFlag, err := configCheck(&cfg) // выводим текущие настройки и если пользователь хочет, может изменить их
		if err != nil {
			fmt.Println("Error::Main::configCheck::", err)
			return
		}
		if changeFlag { // если пользователь решил изменить настройки
			err = configSave(&cfg, configMaster) // запускаем функцию заполнения и сохранения настроек
			if err != nil {
				fmt.Println("Error::Main::configSave::", err)
				return
			}
			cfg, err = configMaster.GetConfig() // считываем записанные настройки
		}
	}

	// Смотрим наличие файлов во входящих
	files, err := ioutil.ReadDir(cfg.InboxPath)
	if err != nil {
		fmt.Println("Error::Main::FileSearch::", err)
		return
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
		if len(addr) > 3 && addr[1:len(addr)-1] == cfg.TargetEmail {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Внимание: Проблема с неверно указанным получателем")

		for i, item := range splice {
			if item != "" {
				fmt.Println("\t", i, item, "!=", cfg.TargetEmail)
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

	err := m.Parser(info, &cfg) // Распределяем данные по структуре письма
	if err != nil {
		fmt.Println("Error::Main::filer::Parser::", err)
		return
	}

	fmt.Println(info.Name())
	action := check_(m) // Производим проверку

	if action {
		err = os.Rename(cfg.InboxPath+"\\"+info.Name(), cfg.SpamPath+"\\"+info.Name())
		if err != nil {
			fmt.Println("Error::Main::filter::FileMove::", err)
			return
		}
	} else {
		err = os.Rename(cfg.InboxPath+"\\"+info.Name(), cfg.FilteredPath+"\\"+info.Name())
		if err != nil {
			fmt.Println("Error::Main::filter::FileMove::", err)
			return
		}
	}
	return
}
