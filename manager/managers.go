package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/Muhammadnumon/alif-bank-console-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
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
	operationsLoop(db,managersOperations,commandOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")

	}
func commandOperationsLoop(db *sql.DB,cmd string)(exit bool){
	switch cmd{
	case "1":
	err:=addClient(db)
		if err != nil {
			log.Printf("can't get all clients:%v",err)
			commandOperationsLoop(db,"1")
			return true
		}
	case "2":
		err := updateBalance(db)
		if err != nil {
			log.Printf("can't add balance: %v", err)
			commandOperationsLoop(db, "2")
			return true
		}
	case "3":
		err := addServices(db)
		if err != nil {
			log.Printf("can't add services: %v", err)
			return true
		}
	case "4":
		err := addBankMachine(db)
		if err != nil {
			log.Printf("can't add bankMachine: %v", err)
			return true
		}
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false

	}
func operationsLoop(db *sql.DB, commands string, loop func(db *sql.DB, cmd string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("can't read input:%v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd))
			exit {
			return
		}
	}
}

///////////////

func addClient(db *sql.DB)(err error){
	fmt.Println("Введите данные клиента")
	fmt.Print("Имя клиента: ")
	reader:=bufio.NewReader(os.Stdin)
	name,err:=reader.ReadString('\n')
	if err!=nil{
		return err
	}
	var login string
	fmt.Print("Логин: ")
	_,err=fmt.Scan(&login)
	if err != nil {
		return err
	}
	var password int64
	fmt.Print("Пароль клиента: ")
	_,err=fmt.Scan(&password)
	if err!=nil{
		return err
	}
	var balance uint64
	fmt.Print("Баланс клиента: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}
	var bankAccount uint64
	fmt.Print("Банковский счёт клиента: ")
	_,err=fmt.Scan(&bankAccount)
	if err != nil {
		return err
	}
	var phoneNumber int64
	fmt.Print("Номер телефон клиента: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}
core.AddClients(core.Client{
	Id:          0,
	Name:        name,
	Login:       login,
	Password:    password,
	BankAccount: bankAccount,
	PhoneNumber: phoneNumber,
	Balance:     balance,
},db)
	if err != nil {
		return err
	}
	fmt.Println("Пользователь успешно добавлен!")
	return nil
}
func addBankMachine(db *sql.DB)(err error){
	fmt.Println("Введите данные банкомата")
	fmt.Print("Имя бокамата: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	var street string
	fmt.Print("Где находится банкомат: ")
	_, err = fmt.Scan(&street)
	if err != nil {
		return err
	}

	err = core.AddBankMachine(core.BankMachine{
		Id:     0,
		Name:   name,
		Street: street,
	}, db)
	if err != nil {
		return err
	}
	fmt.Println("Банкомат успешно добавлен!")
	return nil
}
func addServices(db *sql.DB)(err error){
	fmt.Print("Название услиги: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	var price uint64
	fmt.Print("Стоимость услуги: ")
	_, err = fmt.Scan(&price)
	if err != nil {
		return err
	}
	err = core.AddServices(core.Services{
		Id:    0,
		Name:  name,
		Price: price,
	}, db)
	if err != nil {
		return err
	}
	fmt.Println("Услуга успешно добавлена!")
	return nil
}
func updateBalance(db *sql.DB)(err error){
	fmt.Println("Введите данные клиента")
	var id int64
	fmt.Print("Введите Id клиента: ")
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}
	var name, login string
	var password int64
	var bankAccount uint64
	var phoneNumber int64

	err = core.UpdateBalance(core.Client{
		Id:            id,
		Name:          name,
		Login:         login,
		Password:      password,
		BankAccount: bankAccount,
		PhoneNumber:   phoneNumber,
		Balance:       balance,

	}, db)
	if err != nil {
		return err
	}
	fmt.Println("Баланс клиента успешно добавлен!")
	return nil
}