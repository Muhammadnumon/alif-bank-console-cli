package main

import (
	"database/sql"
	"fmt"
	"github.com/Muhammadnumon/bank-console-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

func main() {
	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer func() {
		log.Print("close db")
		if err := db.Close()
			err != nil {
			log.Fatalf("can't close db: %v", err)
		}
	}()
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init db:%v", err)
	}
	fmt.Println("Добро пожаловать в наше приложение")
	log.Print("start operations loop")
	operationsLoop(db,0, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")
}
func operationsLoop(db *sql.DB, userId int64,commands string, loop func(db *sql.DB, cmd string,userId int64) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("can't read input:%v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd),userId)
			exit {
			return
		}
	}
}
func unauthorizedOperationsLoop(db *sql.DB, cmd string,userId int64) (exit bool) {
	switch cmd {
	case "1":
		userId,_, err := handleLogin(db)
		if err != nil {
			if err==sql.ErrNoRows {
				log.Printf("can't handle login:%v", err)
			}
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")
			return false}
		operationsLoop(db, userId,authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

func authorizedOperationsLoop(db *sql.DB, cmd string,userId int64) (exit bool) {
	switch cmd {
	case "1":
		bankAccounts, err := core.Account(db,userId)
		if err != nil {
			log.Fatalf("can't see bankAccount:%v", err)
			return true
		}
		printBankAccounts(bankAccounts)
	case "2":
		err := transferByPhone(db)
		if err != nil {
			log.Printf("can't transfer by phone:%v", err)
			authorizedOperationsLoop(db, "2",0)
			return true
		}

	case "3":
		err := transferByBankAccount(db)
		if err != nil {
			log.Printf("can't transfer by bankAccount:%v", err)
			authorizedOperationsLoop(db, "3",0)
			return true
		}
	case "4":
		err := payServices(db)
		if err != nil {
			log.Printf("can't pay for services:%v",err)
			authorizedOperationsLoop(db,"4",0)
			return true
		}
		case "5":
		bankMachines, err := core.Machine(db, userId)
		if err != nil {
			log.Fatalf("can't see bankMachine:%v", err)
			return true
		}
		printBankMachines(bankMachines)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

/////////
func printBankMachines(machines []core.BankMachine) {
	for _, machine := range machines {
		fmt.Printf("id:%d,name:%s,street:%s\n",
			machine.Id,
			machine.Name,
			machine.Street)
	}
}
func printBankAccounts(accounts []core.BankAccounts) {
	for _, account := range accounts {
		fmt.Printf("id:%d,name:%s,bankAccount:%d,balance:%d \n",
			account.Id,
			account.Name,
			account.BankAccount,
			account.Balance)
	}
}
func handleLogin(db *sql.DB) (userId int64,ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return 0, false, err
	}
	var password int
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return 0, false, err
	}

	userId,ok, err = core.Login(login, password, db)
	if err != nil {
		return 0,false, err
	}
	return userId,true, nil
}
func payServices(db *sql.DB)(err error){
	var idUser int64
	fmt.Print("Введите свой id:")
	_, err = fmt.Scan(&idUser)
	if err != nil {
		return err
	}
	var idService int64
	fmt.Print("Введите id услуги:")
	_, err = fmt.Scan(&idService)
	if err != nil {
		return err
	}
	var price uint64
	fmt.Print("Введите cтоимость услуги")
	_, err = fmt.Scan(&price)
	if err != nil {
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}
	err=core.PayServicesMinus(core.Client{
		Id:          idUser,
		Balance:     balance,
	},db)
	if err != nil {

		fmt.Println("Извините у ваc недостаточно средств")
		return err
	}
	err = core.PayServicesPlus(idService, price, db)
	if err != nil {
		return err
	}
	fmt.Println("Услуга оплачено!")
	return nil
}
func transferByPhone(db *sql.DB) (err error) {

	var PhoneNumber int64
	fmt.Print("Введите свой номер телефона: ")
	_, err = fmt.Scan(&PhoneNumber)
	if err != nil {
		return err
	}
	var phoneNumber int64
	fmt.Print("Введите номер клиента: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	err = core.TransferMinusByPhoneNumber(core.Client{
		Id:          0,
		Balance:     balance,
		PhoneNumber: PhoneNumber,
	}, db)

	if err != nil {

		fmt.Println("Извините у ваc недостаточно средств")
		return err
	} else {
		if PhoneNumber == phoneNumber {
			fmt.Println("Неверный номер")
			authorizedOperationsLoop(db, "2",0)
		}
	}

	err = core.TransferPlusByPhoneNumber(phoneNumber, balance, db)
	if err != nil {
		return err
	}
	fmt.Println("Счету получаемого успешно отправлено!")
	return nil
}
func transferByBankAccount(db *sql.DB) (err error) {

	var BankAccount uint64
	fmt.Print("Введите номер своего cчета: ")
	_, err = fmt.Scan(&BankAccount)
	if err != nil {
		return err
	}
	var balanceNumber uint64
	fmt.Print("Введите cчет получаемого: ")
	_, err = fmt.Scan(&balanceNumber)
	if err != nil {
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	err = core.TransferMinusByBankAccount(core.Client{
		Id:          0,
		Balance:     balance,
		BankAccount: BankAccount,
	}, db)

	if err != nil {

		fmt.Println("Извините у вас недостаточно средств")
		return err
	} else {
		if BankAccount == balanceNumber {
			fmt.Println("Неверный cчет")
			authorizedOperationsLoop(db, "3",0)
		}
	}

	err = core.TransferPlusByBankAccount(balanceNumber, balance, db)
	if err != nil {
		return err
	}
	fmt.Println("Счет получаемого успешно отправлено!")
	return nil
}
