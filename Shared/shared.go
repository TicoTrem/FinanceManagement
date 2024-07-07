package shared

import (
	"database/sql"
	"log"
	"time"
)

var Database *sql.DB

func SomeFunction() string {
	return "Hello from shared package"
}

type Transaction struct {
	Amount float32
	Date   time.Time
}

// This function will return all of the transactions in the Transactions table
func GetAllTransactions( change this to take in a date range) []Transaction {

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
		err := rows.Scan(&id, &transaction.Amount, &dateString)
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
