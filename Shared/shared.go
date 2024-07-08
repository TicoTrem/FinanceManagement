package shared

import (
	"database/sql"
	"fmt"
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
// if you supply the dBegin and dEnd with nil, it will return all transactions
func GetAllTransactions(dBegin *time.Time, dEnd *time.Time) []Transaction {

	var transactions []Transaction
	// get the first of last month
	today := time.Now()
	if today.Day() != 1 {
		log.Fatal("The monthly task was ran on a day other than the 1st. Please fix this!")
	}

	firstOfLastMonth := today.AddDate(0, -1, 0)

	// this will add a month, the subtract the amount of days, which takes us to the last day of the month
	lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -firstOfLastMonth.Day())
	format := "2006-01-02"
	var rows *sql.Rows
	var err error
	if dBegin == nil || dEnd == nil {
		rows, err = Database.Query(fmt.Sprintf("SELECT * FROM Transactions WHERE date BETWEEN ('%v', '%v')", firstOfLastMonth.Format(format), lastOfLastMonth.Format(format)))
	} else {
		rows, err = Database.Query("SELECT * FROM Transactions")
	}
	if err != nil {
		log.Fatal("Querying all transactions last month failed: " + err.Error())
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
