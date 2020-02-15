package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	core "github.com/shohrukh56/myBank-core"
	"log"
)

const unauthorizedOperations = `
Список доступных операций
1 - вход
2 - список банкоматов
q - Выход из программы
Введите команду`
const authorizedOperations  =`
Список доступных операций
1. Посмотреть список счетов
2. Перевести деньги другому клиенту:
3. Оплатить услугу
q - Выход  в гланое меню
Введите команду
`
const transferOperations  =`
Список доступных операций
1. Перевести деньги по id счета
2. Перевести деньги по номеру телефона:
q - Выход  в гланое меню
Введите команду
`

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open data base %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("can't close data base %v", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		log.Fatalf("can't ping data base %v", err)
	}
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init database %v", err)
	}
	fmt.Println("Добро пожаловать")
	unauthorizedOperationsLoop(db)
}

func unauthorizedOperationsLoop(db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println("\nВведите команду..", unauthorizedOperations)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
		}
		switch cmd {
		case "1":
			fmt.Println("Заполните логин и пароль...")
			handleLogin(db)
		case "2":
			fmt.Println("Список банкоматов...")
			checkCashMachinesList(db)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleLogin(db *sql.DB) {
	var login string
	fmt.Println("Логин: ")
	_, err := fmt.Scan(&login)
	if err != nil {
		fmt.Errorf("can't scan login %w", err)
	}
	var password string
	fmt.Println("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Errorf("can't scan password %w", err)
	}
	dbClientId, name, err := core.Login(db, login, password)

	if err != nil {
		fmt.Println("Не удалось войти в систему! неправильный логин или пароль")
		return
	}
	fmt.Printf("Вы вошли в систему как %s!\n", name)
	authorizedOperationsLoop(dbClientId, db)
}

func authorizedOperationsLoop(dbClientId int, db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println("\nДоступные вам команды...", authorizedOperations)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
			continue
		}
		switch cmd {
		case "1":
			fmt.Println("Список вашых счетов...")
			CheckClientAccounts(db, dbClientId)
		case "2":
			fmt.Println("Перевод денег другому клиенту...")
			handleTransfer(db, dbClientId)
		case "3":
			fmt.Println("Оплатить услугу...")
			handlePayForService(db, dbClientId)
			case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func CheckClientAccounts(db *sql.DB, dbClientId int) {
	list, err := core.ClientAccounts(db, dbClientId)
	if err != nil {
		fmt.Errorf("can't get Bills list! %w", err)
		return
	}
	fmt.Printf("%s\t%s\t%s\n", "id", "Баланс", "Заблокирован")
	for _, value := range list {
		fmt.Printf("%v\t%v\t%v\n", value.Id, value.Balance, value.Locked)
	}
}

func handleTransfer(db *sql.DB, dbClientId int) {
	for ; ; {
		var cmd string
		fmt.Println("\nДоступные вам команды...", transferOperations)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
			continue
		}
		switch cmd {
		case "1":
			fmt.Println("по номеру счёта...")
			handleTransferByAccount(db, dbClientId)
		case "2":
			fmt.Println("по номеру телефона...")
			handleTransferByPhone(db, dbClientId)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleTransferByAccount(db *sql.DB, dbClientId int) {
	var transferAccountId int
	fmt.Print("Введите Id счет пользователя которому хотите осуществить перевод:\nid: ")
	_, err := fmt.Scan(&transferAccountId)
	if err != nil {
		fmt.Printf("can't scan addressee id %v\n", err)
		return
	}
	ok, err, transferAccountBalance := core.CheckAccount(db, transferAccountId)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if !ok {
		fmt.Printf("счет id: %v не существует!", transferAccountId)
		return
	}
	var amount int
	fmt.Print("Введите сумму перевода:\nсумма: ")
	_, err = fmt.Scan(&amount)
	if err != nil {
		fmt.Printf("can't scan amount %v\n", err)
		return
	}
	if amount < 1 {
		fmt.Printf("Сумма перевода долбжна быть больше 0 %v\n", err)
		return
	}
	fmt.Println("доступные вам счета с которых можно осуществить перевод:")
	accounts, err := core.AvailableAccounts(db, dbClientId, amount)
	if err != nil {
		fmt.Printf("Произошла ошибка!!! %v", err)
		return
	}
	fmt.Printf("%s\t%s\n", "id", "Баланс")
	for _, account := range accounts {
		fmt.Printf("%v\t%v\n", account.Id, account.Balance)
	}
	var chosedId int
	fmt.Println("Введите id счета с которого перевести:")
	_, err = fmt.Scan(&chosedId)
	if err != nil {
		fmt.Printf("can't scan account id %v", err)
	}
	for _, value := range accounts {
		if value.Id == chosedId {
			err = core.TransferAccountToAccount(db, value.Id, value.Balance, transferAccountId, transferAccountBalance, amount)
			if err != nil {
				fmt.Printf("не удалось осуществить перевод %v", err)
				return
			}
			fmt.Println("Перевод успешно выполнен.")
			return
		}
	}
	fmt.Printf("Нет такого счета в списке! вы ввели id: %v", chosedId)
}

func handleTransferByPhone(db *sql.DB, dbClietnId int) {
	var phoneAddressOfReceiver string
	fmt.Print("Введите номер телефона пользователя которому хотите осуществить перевод:\nPhone Number: ")
	_, err := fmt.Scan(&phoneAddressOfReceiver)
	if err != nil {
		fmt.Printf("can't scan phone %v\n", err)
		return
	}
	var amount int
	fmt.Print("Введите сумму перевода:\nсумма: ")
	_, err = fmt.Scan(&amount)
	if err != nil {
		fmt.Printf("can't scan amount %v\n", err)
		return
	}
	if amount < 1 {
		fmt.Printf("Сумма перевода долбжна быть больше 0 %v\n", err)
		return
	}
	addressIdOfReceiver, addressBalanceOfReceiver, err := core.FindAccount(db, phoneAddressOfReceiver, amount)

	fmt.Println("доступные вам счета с которых можно осуществить перевод:")
	accounts, err := core.AvailableAccounts(db, dbClietnId, amount)
	if err != nil {
		fmt.Printf("ошибка%v", err)
		return
	}
	fmt.Printf("%s\t%s\n", "id", "Баланс")
	for _, bill := range accounts {
		fmt.Printf("%v\t%v\n", bill.Id, bill.Balance)
	}
	var chosedId int
	fmt.Println("Введите id счета с которого перевести:")
	_, err = fmt.Scan(&chosedId)
	if err != nil {
		fmt.Printf("can't scan bill id %v", err)
	}
	for _, value := range accounts {
		if value.Id == chosedId {
			err = core.TransferAccountToAccount(db, value.Id, value.Balance, addressIdOfReceiver, addressBalanceOfReceiver, amount)
			if err != nil {
				fmt.Printf("не удалось осуществить перевод %v", err)
				return
			}
			fmt.Println("Перевод Успешно выполнен.")
			return
		}
	}
	fmt.Printf("Нет такого счета в списке! вы ввели id: %v", chosedId)
}

func handlePayForService(db *sql.DB, dbClientId int) {
	var service_id int
	fmt.Print("Введите ID услуги: ")
	_, err := fmt.Scan(&service_id)
	if err != nil {
		fmt.Errorf("can't scan service_id %w", err)
		return
	}
	err = core.PayForService(db, service_id, dbClientId)
	if err != nil {
		fmt.Printf("Не удалось оплатить услугу! %v", err)
		return
	}
	fmt.Printf("Услуга : %v была успешно оплачена.\n", service_id)
}

func checkCashMachinesList(db *sql.DB) {
	list, err := core.CashMachinesList(db)
	if err != nil {
		fmt.Errorf("can't make cash machines  list! %w", err)
		return
	}
	fmt.Printf("%s\t%s\t%s\n", "id", "Адрес", "Заблокирован")
	for _, value := range list {
		fmt.Printf("%v\t%s\t%v\n", value.Id, value.Address, value.Locked)
	}
}
var NoSuchAccount = errors.New("нет такого счета")
var NoEnoughMoney = errors.New("недостаточно средств")
var NegativeNumber = errors.New("Отрицательный счет")
var AUserWithThatUserNameAlreadyExists = errors.New("Пользователь с таким логином уже существует")
var AUserWithTheSamePhoneNumberAlreadyExists = errors.New("Пользователь с таким номером телефона уже существует")
var CanNotAdded = errors.New("Таблица уже заполнена!")
var CanNotUnmarshal = errors.New("Can't Unmarshal file")
var CanNotReadFile = errors.New("Can't read file")