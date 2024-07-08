package shared

import (
	"fmt"
	"log"
)

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
func SetEstimatedSpendingMoney(estimatedSpendingMoney float32) {
	_, err := Database.Exec("UPDATE Variables SET estimatedSpendingMoney = ?", estimatedSpendingMoney)
	if err != nil {
		log.Fatal("Failed to update the estimatedSpendingMoney variable: " + err.Error())
	}
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
func SetSpendingMoney(spendingMoney float32) {
	_, err := Database.Exec("UPDATE Variables SET spendingMoney = ?", spendingMoney)
	if err != nil {
		log.Fatal(err)
	}
}

func AddMonthlyExpense(expense MonthlyExpense) {
	_, err := Database.Exec("INSERT INTO MonthlyExpenses (name, amount) VALUES (?, ?);")
	if err != nil {
		log.Fatal("Error inserting expense in to the database" + err.Error())
	}
	fmt.Println("Your monthly expense has successfully been added to the database!")
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

func AddTransaction(transaction Transaction) {
	_, err := Database.Exec("INSERT INTO Transactions (amount, date) VALUES (?, ?);", transaction.Amount, transaction.Date)
	if err != nil {
		log.Fatal("Error inserting transaction in to the database" + err.Error())
	}
	fmt.Println("Your transaction has successfully been added to the database!")
}

// updates the given transaction object in the mysql database. It uses the id in the transaction struct
// to update the correct record, updating the Transaction struct does not update the database
func UpdateTransaction(transaction Transaction) {
	_, err := Database.Exec("UPDATE Transactions SET amount = ? WHERE id = ?", transaction.Amount, transaction.Id)
	if err != nil {
		log.Fatal("Failed to update the transaction: " + err.Error())
	}
}
