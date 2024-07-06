package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

type Transaction struct {
	amount float32
	date   time.Time
}

func main() {

	db, err := sql.Open("mysql", "root:password@/Finance")
	Database = db
	// You started the application before ever running the service

	// how can I make sure the service is currently running?
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Spending money is: %v", calculateSpendingMoney())

	fmt.Println(`Welcome to Finance!\nWhat would you like to do?
						1) Add a transaction
						2) Display all transactions
						3) Change 'Expected' values
						4) View and edit monthly expenses
						5) Add a new goal to save up for`)

	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	for {
		parsedInt, err := strconv.Atoi(response)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		switch parsedInt {
		case 1:
			handleAddTransaction()
		case 2:
			printTransactions(getAllTransactions())
		case 3:
			continue
		default:
			fmt.Println("Invalid input")
			continue
		}
		break
	}

}

func handleAddTransaction() {
	for {
		fmt.Println("What is the amount of the transaction?")
		var amountString string
		_, err := fmt.Scanln(&amountString)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		parsedFloat, err := strconv.ParseFloat(amountString, 32)
		amount := float32(parsedFloat)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		addTransaction(amount, time.Now())
		break
	}

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

// This function will return all of the transactions in the Transactions table
func getAllTransactions() []Transaction {

	var transactions []Transaction
	rows, err := Database.Query("SELECT * FROM Transactions;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

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
	}

	return transactions
}

func printTransactions(transactions []Transaction) {
	for i := 0; i < len(transactions); i++ {
		fmt.Printf("Transaction %v:\nAmount: %v\n Date: %v\n", i+1, transactions[i].amount, transactions[i].date)
	}
}

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() float32 {
	transactions := getAllTransactions()

	var spendingMoney float32 = 0.0
	for i := 0; i < len(transactions); i++ {
		spendingMoney += float32(transactions[i].amount)
	}
	return spendingMoney
}
