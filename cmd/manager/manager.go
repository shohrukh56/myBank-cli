package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	core "github.com/shohrukh56/myBank-core"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const operations = `Список доступных операций:
1. Добавить пользователя
2. Добавить счёт пользователю (тогда сразу в пользователе)
3. Добавить услуги (название)
4. Добавить банкомат
5. Экспорт (форматы json и xml):
6. Импорт того же самого из тех же самых форматов
q. Выйти из приложения
Введите команду`
const exportOperations = `
Что вы хотите экспортировать?
1 - список пользователей
2 - список счетов (с пользователями)
3 - список банкоматов
q - Выход
Введите команду`

func main() {

	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer func() {
		log.Print("close db")
		if err := db.Close(); err != nil {
			log.Fatalf("can't close db: %v", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		log.Fatalf("can't Ping to db: %v", err)
	}
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init db: %v", err)
	}
	fmt.Println("Добро пожаловать в наше приложение")
	log.Print("start operations loop")
	LoopOperation(db)
	log.Print("finish operations loop")
	log.Print("finish application")

}
func LoopOperation(db *sql.DB) {
	for ; ; {
		fmt.Println(operations)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read operation: %v", err)
		}
		switch cmd {
		case "1":
			checkAddClient(db)
		case "2":
			checkAddAccountToClient(db)
		case "3":
			checkAddService(db)
		case "4":
			checkAddСashMachine(db)
		case "5":
			fmt.Println("Export...")
			checkExport(db)
		case "6":
			fmt.Println("Import...")
			checkImport(db)
		case "q":
			return
		default:
			fmt.Println("Неправильная команда: %s\n", cmd)
		}
	}
}
func checkAddClient(db *sql.DB) {
	for ; ; {
		var name string
		var surname string
		var phoneNumber string
		var login string
		var password string
		fmt.Println("Введите Имя")
		_, err := fmt.Scan(&name)
		if err != nil {
			fmt.Errorf("can't read name %w", err)
			fmt.Println("Ошибка. Введите заново")
			return
		}
		fmt.Println("Введите Фамилию")
		_, err = fmt.Scan(&surname)
		if err != nil {
			fmt.Errorf("can't read surname %w", err)
			fmt.Println("Ошибка. Введите заново")
			return
		}
		fmt.Println("Введите номер телефона")
		_, err = fmt.Scan(&phoneNumber)
		if err != nil {
			fmt.Errorf("can't read phone number %w", err)
			fmt.Println("Ошибка. Введите заново")
			return
		}
		fmt.Println("Введите логин")
		_, err = fmt.Scan(&login)
		if err != nil {
			fmt.Errorf("can't read login  %w", err)
			fmt.Println("Ошибка. Введите заново")
			return
		}
		fmt.Println("Введите пароль")
		_, err = fmt.Scan(&password)
		if err != nil {
			fmt.Errorf("can't read password %w", err)
			fmt.Println("Ошибка. Введите заново")
			return
		}
		err = core.AddClient(db, login, password, name, surname, phoneNumber, false)
		if err != nil {
			fmt.Errorf("can't add new client")
			fmt.Println("Несмог добавить нового пользователя в базу данных")
			fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
			var menu string
			_, err = fmt.Scan(&menu)
			if err != nil{
				fmt.Errorf("can't read command")
				fmt.Println("Ошибка. Введите заново")
				return
			}
			return
		} else {
			fmt.Println("Новый пользователь успешно добавлен в базу данных")
			fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
			var menu string
			_, err = fmt.Scan(&menu)
			if err != nil{
				fmt.Errorf("can't read command")
				fmt.Println("Ошибка. Введите заново")
				return
			}
			return
		}
	}
}
func checkAddAccountToClient(db *sql.DB)  {
	var clientID int
    var	balance int
	fmt.Println("Введите id клиента")
	_, err :=fmt.Scan(&clientID)
	if err !=nil{
		fmt.Errorf("can't read client ID ")
		fmt.Println("Ошибка. Введите заново.")
		return
	}
	fmt.Println("Введите баланс клиента")
	_, err =fmt.Scan(&balance)
	if err !=nil{
		fmt.Errorf("can't read balance ")
		fmt.Println("Ошибка. Введите заново.")
		return
	}
	err = core.AddAccountToClient(db, clientID, balance, false)
	if err != nil {
		fmt.Errorf("can't add account to client")
		fmt.Println("Несмог добавить счет пользователю")
		fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
		var menu string
		_, err = fmt.Scan(&menu)
		if err != nil{
			fmt.Errorf("can't read command")
			fmt.Println("Ошибка. Введите заново")
			return
		}
		return
	} else {
		fmt.Println("Новый счет успешно добавлен")
		fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
		var menu string
		_, err = fmt.Scan(&menu)
		if err != nil{
			fmt.Errorf("can't read command")
			fmt.Println("Ошибка. Введите заново")
			return
		}
		return
	}
}
func checkAddService(db *sql.DB)  {
	fmt.Println("Введите название услуги")
	reader := bufio.NewReader(os.Stdin)
	serviceName, err := reader.ReadString('\n')
	if err !=nil{
		fmt.Errorf("can't read serviceName ")
		fmt.Println("Ошибка. Введите заново")
		return
	}
	var servicePrice int
	fmt.Println("Введите цену услуги")
	_, err =fmt.Scan(&servicePrice)
	if err !=nil{
		fmt.Errorf("can't read servicePrice ")
		fmt.Println("Ошибка. Введите заново")
		return
	}
	err = core.AddService(db, serviceName, servicePrice)
	if err != nil {
		fmt.Errorf("can't add account to client")
		fmt.Println("Несмог добавить услугу")
		fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
		var menu string
		_, err = fmt.Scan(&menu)
		if err != nil{
			fmt.Errorf("can't read command")
			fmt.Println("Ошибка. Введите заново")
			return
		}
		return
	} else {
		fmt.Println("Новая услуга успесшно добавлена в базу данных")
		fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
		var menu string
		_, err = fmt.Scan(&menu)
		if err != nil{
			fmt.Errorf("can't read command")
			fmt.Println("Ошибка. Введите заново")
			return
		}
		return
	}
}
func checkAddСashMachine(db *sql.DB){
	fmt.Print("Аддресс банкомата: ")
	reader := bufio.NewReader(os.Stdin)
	address, err := reader.ReadString('\n')
	if err !=nil{
		fmt.Errorf("can't read address ")
		fmt.Println("Ошибка. Введите заново")
		return
	}
	err = core.AddCashMachine(db, address, false)
	if err != nil {
		fmt.Errorf("can't add cash Machine")
		fmt.Println("Несмог добавить банкомат")
		fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
		var menu string
		_, err = fmt.Scan(&menu)
		if err != nil{
			fmt.Errorf("can't read command")
			fmt.Println("Ошибка. Введите заново")
			return
		}
		return
	} else {
		fmt.Println("Новый банкомат успешно добавлен")
		fmt.Println("Введите любую клавишу чтобы выйти в главное меню")
		var menu string
		_, err = fmt.Scan(&menu)
		if err != nil{
			fmt.Errorf("can't read command")
			fmt.Println("Ошибка. Введите заново")
			return
		}
		return
	}
}
func checkExport(db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println(exportOperations)
		_, err2 := fmt.Scan(&cmd)
		if err2!=nil{
			fmt.Println("can't scan cmd")
			continue
		}
		{	switch cmd {
			case "1":
				list, err := core.ClientsList(db)
				if err != nil {
					log.Printf("Can't get clients: %v", err)
					fmt.Println("Не получилосӣ получитӣ список полӣзователей")
					return
				}
				err = exportUsersXmlOrJson(list)
				if err != nil {
					log.Printf("Can't export users: %v", err)
					fmt.Println("Не получилосӣ получитӣ список полӣзователей")
					return
				}
				fmt.Println("Успешный экспорт!")
				return
			case "2":
				list, err := core.AccauntClientList(db)
				if err != nil {
					log.Printf("Can't get Accautslist: %v", err)
					fmt.Println("Не получилосӣ получитӣ список счетов пользователя")
					return
				}
				err = exportAccountsXmlOrJson(list)
				if err != nil {
					log.Printf("Can't export accounts: %v", err)
					fmt.Println("Не получилосӣ получитӣ список Accou")
					return
				}
				fmt.Println("Успешный экспорт!")
				return
			case "3":
				list, err := core.CashMachinesList(db)
				if err != nil {
					log.Printf("Can't read a file: %v", err)
					fmt.Println("Не получилосӣ получитӣ список полӣзователей")
					return
				}
				err = exportCashMachineXmlOrJson(list)
				if err != nil {
					log.Printf("Can't export Cashmachine err: %v", err)
					fmt.Println("Не получилосӣ получитӣ список полӣзователей")
					return
				}
				fmt.Println("Успешный экспорт!")
				return
			case "q":
				return
			default:
				fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
			}
		}
	}
}
func exportUsersXmlOrJson( list []core.ClientList) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("users.json", bytes, 0666)
	if err != nil {
		return err
	}
	bytes, err = xml.Marshal(list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("users.xml", bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func exportCashMachineXmlOrJson( list []core.CashMachine) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("cashmachine.json", bytes, 0666)
	if err != nil {
		return err
	}
	bytes, err = xml.Marshal(list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("cashmachine.xml", bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func exportAccountsXmlOrJson( list []core.ClientAccount) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("accounts.json", bytes, 0666)
	if err != nil {
		return err
	}
	bytes, err = xml.Marshal(list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("accounts.xml", bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func checkImport(db *sql.DB)  {
	var file string
	fmt.Println("Введите название файла полностью")
	fmt.Scan(&file)
	if strings.HasSuffix(file, ".json"){
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("не удалось прочитать файл")
			return
		}
		if strings.HasPrefix(file,"users."){
			list := []core.ClientList{}
			err := json.Unmarshal(bytes, &list)
			if err != nil {
				fmt.Println("Не удалось конвертировать")
				return
			}
			for _, clientList := range list {
				fmt.Println(clientList)
			}
		}
		if strings.HasPrefix(file,"accounts."){
			list := []core.ClientAccount{}
			err := json.Unmarshal(bytes, &list)
			if err != nil {
				fmt.Println("Не удалось конвертировать")
				return
			}
			for _, clientList := range list {
				fmt.Println(clientList)
			}
		}
		if strings.HasPrefix(file,"cashmachine."){
			list := []core.CashMachine{}
			err := json.Unmarshal(bytes, &list)
			if err != nil {
				fmt.Println("Не удалось конвертировать")
				return
			}
			for _, clientList := range list {
				fmt.Println(clientList)
			}
		}

	}else 	if strings.HasSuffix(file, ".xml"){
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("не удалось прочитать файл")
			return
		}
		if strings.HasPrefix(file,"users."){
			list := []core.ClientList{}
			err := xml.Unmarshal(bytes, &list)
			if err != nil {
				fmt.Println("Не удалось конвертировать")
				return
			}
			for _, clientList := range list {
				fmt.Println(clientList)
			}
		}
		if strings.HasPrefix(file,"accounts."){
			list := []core.ClientAccount{}
			err := xml.Unmarshal(bytes, &list)
			if err != nil {
				fmt.Println("Не удалось конвертировать")
				return
			}
			for _, clientList := range list {
				fmt.Println(clientList)
			}
		}
		if strings.HasPrefix(file,"cashmachine."){
			list := []core.CashMachine{}
			err := xml.Unmarshal(bytes, &list)
			if err != nil {
				fmt.Println("Не удалось конвертировать")
				return
			}
			for _, clientList := range list {
				fmt.Println(clientList)
			}
		}
	}else {
		fmt.Println("Ошибка формата")
	}
}