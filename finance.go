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

var Database *sql.DB

type Transaction struct {
	amount float32
	date   time.Time
}

func StartFinance() {
	go queueMonthlyTask()

	// create database if it doesn't already exist
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

	// Create the database for real
	Database, err = sql.Open("mysql", "root:password@/Finance")
	if err != nil {
		log.Fatal(err)
	}
	defer Database.Close()
	fmt.Println("3")

	createTables()

	// get the starting spending money (intensive operation)
	calculateSpendingMoney()
	addTransaction(53.2, time.Now())
	var transactions []Transaction = getAllTransactions()

	printTransactions(transactions)

	fmt.Println(db.Ping())
}

func printTransactions(transactions []Transaction) {
	for i := 0; i < len(transactions); i++ {
		fmt.Printf("Transaction %v:\nAmount: %v\n Date: %v\n", i+1, transactions[i].amount, transactions[i].date)
	}
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

func addTransaction(amount float32, date time.Time) {
	transaction := Transaction{amount: amount, date: date}
	query, err := Database.Prepare("INSERT INTO Transactions (amount, date) VALUES (?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	result, err := query.Exec(transaction.amount, transaction.date)
	if err != nil {
		log.Fatal(err)
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There were %v rows inserted into the Transactions table\n", numRows)
}

// creates the tables needed for the application if they are not created already
// also populates the Variables table with a row containing all 0.0 if there is not already a row
func createTables() {
	_, err := Database.Exec(`CREATE TABLE IF NOT EXISTS Transactions (
		id INT AUTO_INCREMENT,
		amount FLOAT(16,2) NOT NULL,
		date DATETIME NOT NULL,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS Variables (
	spendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedSpendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedIncome FLOAT(16,2) DEFAULT 0.0,
	totalChangeThisMonth FLOAT(16,2) DEFAULT 0.0);`)
	if err != nil {
		log.Fatal(err)
	}
	row := Database.QueryRow(`SELECT COUNT(*) FROM Variables`)
	var count int
	err = row.Scan(&count)

	if err != nil {
		log.Fatal(err)
	}
	// there are no rows in the table
	if count == 0 {
		Database.Exec(`INSERT INTO Variables () VALUES ()`)
	}

	// numRows, err := result.RowsAffected()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // the table was just created, initialize values to defaults
	// if numRows > 0 {

	// }
}

// This function will return all of the transactions in the Transactions table
// func getAllTransactions() []Transaction {

// 	var transactions []Transaction
// 	rows, err := Database.Query("SELECT * FROM Transactions;")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var transaction Transaction
// 		var id int
// 		var dateString string
// 		// sql date is wanting to return a string
// 		err := rows.Scan(&id, &transaction.amount, &dateString)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		parsedDate, err := time.Parse("2006-01-02 15:04:05", dateString)
// 		if err != nil {
// 			log.Fatal("Failed to parse SQL string into a time object:", err)
// 		}
// 		transaction.date = parsedDate
// 		transactions = append(transactions, transaction)
// 	}

// 	return transactions
// }
