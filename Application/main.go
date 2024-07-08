package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ticotrem/shared"
)

var Database *sql.DB

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
			printTransactions(shared.GetAllTransactions())
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
	transaction := shared.Transaction{Amount: amount, Date: date}
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

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() float32 {
	transactions := shared.GetAllTransactions()

	var spendingMoney float32 = 0.0
	for i := 0; i < len(transactions); i++ {
		spendingMoney += float32(transactions[i].Amount)
	}
	return spendingMoney
}

func printTransactions(transactions []shared.Transaction) {
	for i := 0; i < len(transactions); i++ {
		fmt.Printf("Transaction %v:\nAmount: %v\n Date: %v\n", i+1, transactions[i].Amount, transactions[i].Date)
	}
}
