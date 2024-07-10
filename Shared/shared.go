package shared

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var Database *sql.DB

// This function will setup the database and create the tables if they don't exist
func SetupDatabase() {

	db, err := sql.Open("mysql", "root:password@/")
	if err != nil {
		log.Fatal(err)
	}

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
    	description VARCHAR(100) NOT NULL,
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

	// consider changing Goals table to store months left instead of dateComplete
	// TODO: Change goals to just keep track of months instead of a date
	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS Goals (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		amount FLOAT(16,2) NOT NULL,
    	amountSaved float(16,2) NOT NULL,
    	dateComplete DATE NOT NULL,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS Variables (
	spendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedSpendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedIncome FLOAT(16,2) DEFAULT 0.0,
    emergencyMax float(16,2) DEFAULT 0.0,
    emergencyAmount float(16,2) DEFAULT 0.0);`)
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
	var differenceOfExpenses float32 = 0.0
	for i := 0; i < len(expenses); i++ {
		differenceOfExpenses -= expenses[i].Amount
	}
	return differenceOfExpenses
}

func MonthlyTask() {

	// add the monthly transactions for that month and set the date as the last day of the previous month
	// this is so you can make changes to the prices of the monthly expenses and have it reflected in that months transactions
	// TODO: make sure the estimated spending money is updated when the user alters a monthly expense value, it is assumed to be active THAT MONTH
	// so it must be updated.

	// add expenses as transactions
	expenses := GetAllMonthlyExpensesStructs()
	for i := 0; i < len(expenses); i++ {
		AddTransaction(&Transaction{Amount: expenses[i].Amount, Date: time.Now().AddDate(0, 0, -1), Description: fmt.Sprintf("Expenses: %v monthly payment", expenses[i].Name)})
	}

	// add goals as transactions
	goals := GetAllGoalStructs()
	for i := 0; i < len(goals); i++ {
		goals[i].SaveMonthlyAmount()
	}

	// the above 2 are in the transactions so we can now calculate how much
	// our money went up or down this month in total (no estimated values)
	netTransactionChange := calculateNetTransactionChange()

	spendingMoney := GetSpendingMoney() + netTransactionChange

	// update it for last month, we later updated EstimatedSpendingMoney for THIS month.
	SetSpendingMoney(spendingMoney)

	// Set the estimated spending money value to the spending money, with next months predicted outcome
	// and deducting the set in stone monthly expenses. The expenses should be automatically registered as transactions because
	// otherwise you would not be able to lower the spending money when you make purchases
	SetEstimatedSpendingMoney(spendingMoney + GetExpectedMonthlyIncome() - GetMonthlyExpenses())

	// The emergency fund takes half of the netTransaction change if it is positive.
	emergencyAmount, emergencyMax := GetEmergencyData()

	UpdateMaxEmergencyFund()
	if emergencyAmount < emergencyMax && netTransactionChange > 0 {
		// creates a transaction, lowering the spending money
		IncreaseEmergencyFund(netTransactionChange / 2)
	}

}

// This will calculate the net transaction change (which includes income and expenses)
// Then this will change estimated spending money to this value
func calculateNetTransactionChange() float32 {

	// get the first of last month
	today := time.Now()
	// TODO: UNCOMMENT THIS WHEN DONE, ITS GONE FOR TESTING
	//if today.Day() != 1 {
	//	log.Fatal("The monthly task was ran on a day other than the 1st. Please fix this!")
	//}

	firstOfLastMonth := today.AddDate(0, -1, 0)

	// this will add a month, the subtract the amount of days, which takes us to the last day of the month
	lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -firstOfLastMonth.Day())

	var netTransactionChange float32 = 0.0

	lastMonthTransactions := GetAllTransactions(&firstOfLastMonth, &lastOfLastMonth)

	for i := 0; i < len(lastMonthTransactions); i++ {
		netTransactionChange += lastMonthTransactions[i].Amount
	}

	return netTransactionChange
}
