package shared

import (
	"database/sql"
	"fmt"
	"log"
)

var Database *sql.DB

// This function will setup the database and create the tables if they don't exist
func SetupDatabase() {

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
	db.Close()

	// Create the database object for real
	Database, err = sql.Open("mysql", "root:password@/Finance")
	if err != nil {
		log.Fatal(err)
	}
	defer Database.Close()

	createTables()

	// get the starting spending money (intensive operation)

	if Database.Ping() != nil {
		log.Fatal("Failed to ping database")
	}
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

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS MonthlyExpenses (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		amount FLOAT(16,2) NOT NULL,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS Goals (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		amount FLOAT(16,2) NOT NULL,
    	dateComplete DATE NOT NULL,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS Variables (
	spendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedSpendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedIncome FLOAT(16,2) DEFAULT 0.0;`)
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

}

func GetMonthlyExpenses() float32 {
	expenses := GetAllMonthlyExpensesStructs()
	var sumOfExpenses float32 = 0.0
	for i := 0; i < len(expenses); i++ {
		sumOfExpenses += expenses[i].Amount
	}
	return sumOfExpenses
}
