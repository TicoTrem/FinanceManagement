package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ticotrem/finance/shared/utils"
)

type Transaction struct {
	Id          int
	Amount      float32
	Date        time.Time
	Description string
}

func AddTransaction(transaction *Transaction) {
	_, err := Database.Exec("INSERT INTO Transactions (amount, date, description) VALUES (?, ?, ?);", transaction.Amount, transaction.Date, transaction.Description)
	if err != nil {
		log.Fatal("Error inserting transaction in to the Database" + err.Error())
	}
}

func (transaction *Transaction) Delete() {
	_, err := Database.Exec("DELETE FROM Transactions WHERE id = ?;", transaction.Id)
	if err != nil {
		log.Fatal("Error deleting transaction %v in the Database"+err.Error(), transaction.Id)
	}
	if transaction.Amount < 0 {
		SetEstimatedSpendingMoney(GetEstimatedSpendingMoney() - transaction.Amount)
	}

	fmt.Println("Your transaction has successfully been deleted from the Database!")
	*transaction = Transaction{} // make it zero value since it is now deleted from Database
}

// IsFromCurrentMonth checks if the transaction date is from the current month
// It returns true if the transaction is from the current month
func (t *Transaction) IsFromCurrentMonth() bool {
	currentTime := utils.CurrentTime()

	// Check if the transaction year and month match the current year and month
	return t.Date.Year() == currentTime.Year() &&
		t.Date.Month() == currentTime.Month()
}

// updates the given transaction object in the mysql Database. It uses the id in the transaction struct
// to update the correct record, updating the Transaction struct does not update the Database
func (transaction *Transaction) Update() {

	_, err := Database.Exec("UPDATE Transactions SET amount = ? WHERE id = ?;", transaction.Amount, transaction.Id)
	if err != nil {
		log.Fatal("Failed to update the transaction: " + err.Error())
	}
}

// This function will return all of the transactions in the Transactions table
// if you supply the dBegin and dEnd with nil, it will return all transactions
func GetAllTransactions(dBegin *time.Time, dEnd *time.Time) []Transaction {

	var transactions []Transaction

	format := "2006-01-02"
	var rows *sql.Rows
	var err error
	if dBegin != nil && dEnd != nil {
		rows, err = Database.Query("SELECT * FROM Transactions WHERE date BETWEEN ? AND ? ORDER BY date;", dBegin.Format(format), dEnd.Format(format))
	} else {
		rows, err = Database.Query("SELECT * FROM Transactions ORDER BY date;")
	}
	if err != nil {
		log.Fatal("Querying transactions failed: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {

		var transaction Transaction
		var dateString string
		// sql date is wanting to return a string
		err := rows.Scan(&transaction.Id, &transaction.Amount, &dateString, &transaction.Description)
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
