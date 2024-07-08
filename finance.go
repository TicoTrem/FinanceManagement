package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var spendingMoney float32 = 0
var estimatedSpendingMoney float32 = 0
var estimatedIncome float32 = 0
var totalChangeThisMonth float32 = 0

var Database sql.DB

func StartFinance() {
	go queueMonthlyTask()
	db, err := sql.Open("mysql", "root:password@/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("1")
	fmt.Println(db.Ping())

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS Finance")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2")
	db.Close()

	db, err = sql.Open("mysql", "root:password@/Finance")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("3")

	createTables()

	// get the starting spending money (intensive operation)
	calculateSpendingMoney()

	fmt.Println(db.Ping())
}

func queueMonthlyTask() {
	for {
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		duration := firstOfNextMonth.Sub(now)

		timer := time.NewTimer(duration)

		// blocks execution until the timer expires (the start of the month)
		<-timer.C
		monthlyTask()
	}
}

func monthlyTask() {
	calculateMonthSpendingMoney()
}

// This will calculate the net transaction change (which includes income and expenses)
// Then this will change estimated spending money to this value
func calculateMonthSpendingMoney() float32 {
	return 5.0
}

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() {
	fmt.Println("calculate spending money called")
}

func createTables() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Transactions (
		id INT AUTO_INCREMENT,
		amount INT NOT NULL,
		date DATE NOT NULL,
		PRIMARY KEY(id));`)
}
