package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

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
	// Database.Query()
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
