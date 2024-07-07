package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)



func StartFinance() {

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

	// Create the database object for real
	Database, err = sql.Open("mysql", "root:password@/Finance")
	if err != nil {
		log.Fatal(err)
	}
	defer Database.Close()
	fmt.Println("3")

	createTables()

	// get the starting spending money (intensive operation)

	fmt.Println(db.Ping())

	go queueMonthlyTask()
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
	// get the first of last month
	today := time.Now()
	if today.Day() != 1 {
		log.Fatal("The monthly task was ran on a day other than the 1st. Please fix this!")
	}

	firstOfLastMonth := today.AddDate(0, -1, 0)

	// this will add a month, the subtract the amount of days, which takes us to the last day of the month
	lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -firstOfLastMonth.Day())
	format := "2006-01-02"
	rows, err := Database.Query(fmt.Sprintf("SELECT * FROM Transactions WHERE date BETWEEN ('%v', '%v')", firstOfLastMonth.Format(format), lastOfLastMonth.Format(format)))
	if err != nil {
		log.Fatal("Querying all transactions last month failed: " + err.Error())
	}

	netTransactionChange := 0

	for rows.Next() {
		var transaction Transaction
		var id int
		var dateString string
		// sql date is wanting to return a string
		err := rows.Scan(&id, &transaction.amount, &dateString)
		if err != nil {
			log.Fatal(err)
		}
		parsedDate, err := time.Parse("2006-01-02 15:04:05", dateString)
		if err != nil {
			log.Fatal("Failed to parse SQL string into a time object:", err)
		}
		transaction.date = parsedDate
		transactions = append(transactions, transaction)

	return 15.9
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
