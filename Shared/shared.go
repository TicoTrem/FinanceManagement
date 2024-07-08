package shared

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var Database *sql.DB

type Transaction struct {
	Id     int
	Amount float32
	Date   time.Time
}

type MonthlyExpense struct {
	Id     int
	Name   string
	Amount float32
}

// This function will return all of the transactions in the Transactions table
// if you supply the dBegin and dEnd with nil, it will return all transactions
func GetAllTransactions(dBegin *time.Time, dEnd *time.Time) []Transaction {

	var transactions []Transaction

	format := "2006-01-02"
	var rows *sql.Rows
	var err error
	if dBegin != nil && dEnd != nil {
		rows, err = Database.Query(fmt.Sprintf("SELECT * FROM Transactions WHERE date BETWEEN ('%v', '%v')", dBegin.Format(format), dEnd.Format(format)))
	} else {
		rows, err = Database.Query("SELECT * FROM Transactions")
	}
	if err != nil {
		log.Fatal("Querying all transactions last month failed: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {

		var transaction Transaction
		var dateString string
		// sql date is wanting to return a string
		err := rows.Scan(&transaction.Id, &transaction.Amount, &dateString)
		if err != nil {
			log.Fatal(err)
		}
		parsedDate, err := time.Parse("2006-01-02 15:04:05", dateString)
		if err != nil {
			log.Fatal("Failed to parse SQL string into a time object:", err)
		}
		transaction.Date = parsedDate
		transactions = append(transactions, transaction)
	}

	return transactions
}

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

func AddTransaction(amount float32, date time.Time) {
	transaction := Transaction{Amount: amount, Date: date}
	query, err := Database.Prepare("INSERT INTO Transactions (amount, date) VALUES (?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	result, err := query.Exec(transaction.Amount, transaction.Date)
	if err != nil {
		log.Fatal(err)
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There were %v rows inserted into the Transactions table\n", numRows)
}

// This function will return the expectedIncome variable from the Variables table
func GetSpendingMoney() float32 {
	row := Database.QueryRow("SELECT spendingMoney FROM Variables")
	var spendingMoney float32
	err := row.Scan(&spendingMoney)
	if err != nil {
		log.Fatal(err)
	}
	return spendingMoney
}

// This function will return the estimatedExpectedIncome variable from the Variables table
func GetEstimatedSpendingMoney() float32 {
	row := Database.QueryRow("SELECT estimatedSpendingMoney FROM Variables")
	var estimatedSpendingMoney float32
	err := row.Scan(&estimatedSpendingMoney)
	if err != nil {
		log.Fatal(err)
	}
	return estimatedSpendingMoney
}

func SetSpendingMoney(spendingMoney float32) {
	_, err := Database.Exec("UPDATE Variables SET spendingMoney = ?", spendingMoney)
	if err != nil {
		log.Fatal(err)
	}
}

func SetEstimatedSpendingMoney(estimatedSpendingMoney float32) {
	_, err := Database.Exec("UPDATE Variables SET estimatedSpendingMoney = ?", estimatedSpendingMoney)
	if err != nil {
		log.Fatal("Failed to update the estimatedSpendingMoney variable: " + err.Error())
	}
}

func GetExpectedMonthlyIncome() float32 {
	row := Database.QueryRow("SELECT estimatedIncome FROM Variables")
	var estimatedIncome float32
	row.Scan(&estimatedIncome)
	return estimatedIncome
}

func SetExpectedMonthlyIncome(expectedMonthlyIncome float32) {
	_, err := Database.Exec("UPDATE Variables SET expectedMonthlyIncome = ?", expectedMonthlyIncome)
	if err != nil {
		log.Fatal("Failed to update the expectedMonthlyIncome variable: " + err.Error())
	}
}

func GetAllMonthlyExpensesStructs() []MonthlyExpense {
	rows, err := Database.Query("SELECT * FROM MonthlyExpenses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var monthlyExpenses []MonthlyExpense
	var id int

	for rows.Next() {

		var monthlyExpense MonthlyExpense

		// sql date is wanting to return a string
		err := rows.Scan(&id, &monthlyExpense.Name, &monthlyExpense.Amount)
		if err != nil {
			log.Fatal(err)
		}

		monthlyExpenses = append(monthlyExpenses, monthlyExpense)
	}

	return monthlyExpenses
}

func GetMonthlyExpenses() float32 {
	expenses := GetAllMonthlyExpensesStructs()
	var sumOfExpenses float32 = 0.0
	for i := 0; i < len(expenses); i++ {
		sumOfExpenses += expenses[i].Amount
	}
	return sumOfExpenses
}
