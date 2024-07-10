package shared

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// TODO: store times in UTC time, then convert to show the user their local time

func GetExpectedMonthlyIncome() float32 {
	row := Database.QueryRow("SELECT estimatedIncome FROM Variables")
	var estimatedIncome float32
	row.Scan(&estimatedIncome)
	return estimatedIncome
}
func SetEstimatedMonthlyIncome(estimatedMonthlyIncome float32) {
	_, err := Database.Exec("UPDATE Variables SET estimatedIncome = ?", estimatedMonthlyIncome)
	if err != nil {
		log.Fatal("Failed to update the estimatedMonthlyIncome variable: " + err.Error())
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

// This function will return the spendingMoney variable from the Variables table
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
	//_, err := Database.Exec("UPDATE Variables SET spendingMoney = ?", spendingMoney)
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func AddMonthlyExpense(expense MonthlyExpense) {
	_, err := Database.Exec(fmt.Sprintf("INSERT INTO MonthlyExpenses (name, amount) VALUES ('%v', %v);", expense.Name, expense.Amount))
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

	for rows.Next() {

		var monthlyExpense MonthlyExpense

		// sql date is wanting to return a string
		err := rows.Scan(&monthlyExpense.Id, &monthlyExpense.Name, &monthlyExpense.Amount)
		if err != nil {
			log.Fatal(err)
		}

		monthlyExpenses = append(monthlyExpenses, monthlyExpense)
	}

	return monthlyExpenses
}

func AddTransaction(transaction *Transaction) {
	_, err := Database.Exec("INSERT INTO Transactions (amount, date) VALUES (?, ?);", transaction.Amount, transaction.Date)
	if err != nil {
		log.Fatal("Error inserting transaction in to the database" + err.Error())
	}
	fmt.Println("Your transaction has successfully been added to the database!")
}

func DeleteTransaction(transaction *Transaction) {
	_, err := Database.Exec("DELETE FROM Transactions WHERE id = ?;", transaction.Id)
	if err != nil {
		log.Fatal("Error deleting transaction in the database" + err.Error())
	}
	fmt.Println("Your transaction has successfully been deleted from the database!")
}

// updates the given transaction object in the mysql database. It uses the id in the transaction struct
// to update the correct record, updating the Transaction struct does not update the database
func UpdateTransaction(transaction *Transaction) {
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
		rows, err = Database.Query(fmt.Sprintf("SELECT * FROM Transactions WHERE date BETWEEN ('%v', '%v') ORDER BY date", dBegin.Format(format), dEnd.Format(format)))
	} else {
		rows, err = Database.Query("SELECT * FROM Transactions ORDER BY date")
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

func AddGoal(goal *Goal) {
	_, err := Database.Exec("INSERT INTO Goals (name, amount, amountSaved, dateComplete) VALUES ('?', ?, ?, '?');", goal.Name, goal.Amount, goal.AmountSaved, goal.DateComplete)
	if err != nil {
		log.Fatal("Error inserting goal in to the database:" + err.Error())
	}
}

func GetAllGoals() []Goal {
	var goals []Goal
	rows, err := Database.Query("SELECT * FROM Goals;")
	if err != nil {
		log.Fatal("Querying all goals failed" + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var goal Goal
		var dateString string
		err = rows.Scan(&goal.Id, &goal.Name, &goal.Amount, &dateString)
		if err != nil {
			log.Fatal("Failed to scan goal into goal struct:" + err.Error())
		}
		parsedDate, err := time.Parse("2006-01-02 15:04:05", dateString)
		if err != nil {
			log.Fatal("Failed to parse SQL string into a time object:", err)
		}
		goal.DateComplete = parsedDate
		// calculate and assign the amount per month attribute
		months := getMonthsToComplete(&goal)
		goal.AmountPerMonth = goal.Amount / float32(months)
		goals = append(goals, goal)
	}
	return goals
}

func getMonthsToComplete(goal *Goal) int {
	months := (time.Now().Year() - goal.DateComplete.Year()) * 12
	months += int(time.Now().Month() - goal.DateComplete.Month())
	return months

}

func (goal *Goal) UpdateGoal() {
	_, err := Database.Exec("UPDATE Goals SET amountSaved = ? WHERE id = ?;", goal.Id, goal.AmountPerMonth)
	if err != nil {
		log.Fatal("Error updating goal in database: " + err.Error())
	}
}

func GetEmergencyData() (max float32, amount float32) {
	row := Database.QueryRow("SELECT emergencyMax, emergencyAmount FROM Variables")
	var emergencyMax float32
	var emergencyAmount float32
	err := row.Scan(&emergencyMax, &emergencyAmount)
	if err != nil {
		log.Fatal("Failed to scan emergency values into variables:" + err.Error())
	}
	return emergencyAmount, emergencyMax
}

func IncreaseEmergencyAmount(amount float32) {
	emergencyAmount, _ := GetEmergencyData()
	_, err := Database.Exec("UPDATE Variables SET emergencyAmount = ?;", emergencyAmount+amount)
	if err != nil {
		log.Fatal("Error updating emergency values into variables table:" + err.Error())
	}
}
