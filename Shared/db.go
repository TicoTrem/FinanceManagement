package shared

import (
	"database/sql"
	"fmt"
	"log"
	"math"
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
func (expense *MonthlyExpense) UpdateExpenseName(name string) {
	_, err := Database.Exec("UPDATE MonthlyExpenses SET name = ? WHERE id = ?;", name, expense.Id)
	if err != nil {
		log.Fatal("Failed to update the expense name: " + err.Error())
	}
}
func (expense *MonthlyExpense) UpdateExpenseAmount(amount float32) {
	_, err := Database.Exec("UPDATE MonthlyExpenses SET amount = ? WHERE id = ?;", amount, expense.Id)
	if err != nil {
		log.Fatal("Failed to update the expense amount: " + err.Error())
	}
}

// assumes the payment was made that month already
func (expense *MonthlyExpense) Delete() {
	_, err := Database.Exec("DELETE FROM MonthlyExpenses WHERE id = ?;", expense.Id)
	if err != nil {
		log.Fatal("Failed to delete the expense from the database: " + err.Error())
	}
	*expense = MonthlyExpense{} // make it zero value since it is now deleted from database
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
	_, err := Database.Exec("INSERT INTO Transactions (amount, date, description) VALUES (?, ?, ?);", transaction.Amount, transaction.Date, transaction.Description)
	if err != nil {
		log.Fatal("Error inserting transaction in to the database" + err.Error())
	}
}

func (transaction *Transaction) Delete() {
	_, err := Database.Exec("DELETE FROM Transactions WHERE id = ?;", transaction.Id)
	if err != nil {
		log.Fatal("Error deleting transaction in the database" + err.Error())
	}
	fmt.Println("Your transaction has successfully been deleted from the database!")
	*transaction = Transaction{} // make it zero value since it is now deleted from database
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
		rows, err = Database.Query("SELECT * FROM Transactions WHERE date BETWEEN ? AND ? ORDER BY date;", dBegin.Format(format), dEnd.Format(format))
	} else {
		rows, err = Database.Query("SELECT * FROM Transactions ORDER BY date;")
	}
	if err != nil {
		log.Fatal("Querying all transactions last month failed: " + err.Error())
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

func AddGoal(goal *Goal) {
	_, err := Database.Exec("INSERT INTO Goals (name, amount, amountSaved, dateComplete) VALUES (?, ?, ?, ?);", goal.Name, goal.Amount, goal.AmountSaved, goal.DateComplete)
	if err != nil {
		log.Fatal("Error inserting goal in to the database:" + err.Error())
	}
}

func GetAllGoalStructs() []Goal {
	var goals []Goal
	rows, err := Database.Query("SELECT * FROM Goals;")
	if err != nil {
		log.Fatal("Querying all goals failed" + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var goal Goal
		var dateString string
		err = rows.Scan(&goal.Id, &goal.Name, &goal.Amount, &goal.AmountSaved, &dateString)
		if err != nil {
			log.Fatal("Failed to scan goal into goal struct:" + err.Error())
		}
		parsedDate, err := time.Parse(time.DateOnly, dateString)
		if err != nil {
			log.Fatal("Failed to parse SQL string into a time object:", err)
		}
		goal.DateComplete = parsedDate
		// calculate and assign the amount per month attribute
		goal.PopulateAmountPerMonth()

		goals = append(goals, goal)
	}
	return goals
}

// populates the AmountPerMonth attribute of the struct, assuming it has Amount and DateComplete specified
func (goal *Goal) PopulateAmountPerMonth() {
	zeroTime := time.Time{}
	if goal.Amount == 0 || goal.DateComplete == zeroTime {
		log.Fatal("You cannot populate the amount per month without having the Amount and DateComplete fields")
	}
	months := goal.GetMonthsToComplete()
	if months <= 0 {
		goal.AmountPerMonth = float32(0)
	} else {
		goal.AmountPerMonth = (goal.Amount - goal.AmountSaved) / float32(months)
	}
}

// populates the DateComplete attribute of the struct, assuming it has AmountPerMonth and Amount specified
func (goal *Goal) PopulateDateComplete() {
	// round up, so if it takes 3.1 months, we will give it 4 months (as they can't afford it in the previous month
	var months int = int(math.Ceil(float64(goal.Amount / goal.AmountPerMonth)))
	now := time.Now()
	goal.DateComplete = time.Date(now.Year(), now.Month()+time.Month(months), 1, 0, 0, 0, 0, now.Location())
}

// too confusing
//
//	func (goal *Goal) GetMonthsToComplete() int {
//		years := goal.DateComplete.Year() - time.Now().Year()
//		months := int(time.Now().Month() - goal.DateComplete.Month())
//
//		totalMonths := (years * 12) + months
//		return totalMonths
//	}

// TODO: understand claude 3.5 code
func (goal *Goal) GetMonthsToComplete() int {
	now := time.Now()
	if now.After(goal.DateComplete) {
		return 0 // Goal is already complete
	}

	years := goal.DateComplete.Year() - now.Year()
	months := int(goal.DateComplete.Month() - now.Month())

	totalMonths := years*12 + months

	// Adjust for day of month
	if goal.DateComplete.Day() < now.Day() {
		totalMonths--
	}

	return totalMonths
}
func (goal *Goal) SaveMonthlyAmount() {
	AddTransaction(&Transaction{Amount: goal.Amount, Description: fmt.Sprintf("Goal: %v monthly savings", goal.Name)})
	_, err := Database.Exec("UPDATE Goals SET amountSaved = ? WHERE id = ?;", goal.AmountSaved+goal.AmountPerMonth, goal.Id)
	if err != nil {
		log.Fatal("Error updating goal in database: " + err.Error())
	}
}

func (goal *Goal) UpdateGoalName(name string) {
	_, err := Database.Exec("UPDATE Goals SET name = ? WHERE id = ?;", name, goal.Id)
	if err != nil {
		log.Fatal("Failed to update the goal name: " + err.Error())
	}
}

// UpdateGoalAmount updates the goal amount in the database, and either increases monthly payment to compensate or increases
// the time to complete the goal as specified from the payMoreMonthly boolean
func (goal *Goal) UpdateGoalAmount(amount float32, payMoreMonthly bool) {
	goal.Amount = amount
	var err error
	// adjust the monthly amount
	if payMoreMonthly {
		// this will calculate and apply a new value since goal.Amount changed
		goal.PopulateAmountPerMonth()
		// we do not need to update the AmountPerMonth in the database as its only something calculated
		_, err = Database.Exec("UPDATE MonthlyExpenses SET amount = ? WHERE id = ?;", amount, goal.Id)
	} else { // adjust the time to completion
		// this will calculate and apply a new value since goal.Amount changed
		goal.PopulateDateComplete()
		_, err = Database.Exec("UPDATE MonthlyExpenses SET amount = ? AND dateComplete = ? WHERE id = ?;", amount, goal.DateComplete, goal.Id)
	}
	if err != nil {
		log.Fatal("Failed to update the expense amount: " + err.Error())
	}
}
func (goal *Goal) UpdateGoalDate(date time.Time) {
	_, err := Database.Exec("UPDATE Goals SET dateComplete = ? WHERE id = ?;", date, goal.Id)
	if err != nil {
		log.Fatal("Failed to update the goal date in database: " + err.Error())
	}
	// update the amount per month with new completion date
	goal.PopulateAmountPerMonth()
}
func (goal *Goal) UpdateGoalMonthly(amountPerMonth float32) {
	goal.AmountPerMonth = amountPerMonth
	// updates the date it will be completed by based on the amount per month and the unchanged amount
	goal.PopulateDateComplete()
	// update date in database
	_, err := Database.Exec("UPDATE Goals SET dateComplete = ? WHERE id = ?;", goal.DateComplete, goal.Id)
	if err != nil {
		log.Fatal("Failed to update the dateComplete based on the new amountPerMonth in database: " + err.Error())
	}
}
func (goal *Goal) DeleteGoal() {
	_, err := Database.Exec("DELETE FROM Goals WHERE id = ?;", goal.Id)
	if err != nil {
		log.Fatal("Error deleting goal from the database: " + err.Error())
	}
	*goal = Goal{} // make it zero value since it is now deleted from database
}
func (goal *Goal) Contribute(amountToContribute float32) {
	_, err := Database.Exec("UPDATE Goals SET amountSaved = ? WHERE id = ?;", goal.AmountSaved+amountToContribute, goal.Id)
	if err != nil {
		log.Fatal("Error contributing to goal in database: " + err.Error())
	}
}

func GetEmergencyData() (max float32, amount float32) {
	row := Database.QueryRow("SELECT emergencyMax, emergencyAmount FROM Variables")
	var emergencyMax float32
	var emergencyAmount float32
	err := row.Scan(&emergencyMax, &emergencyAmount)
	if err != nil {
		log.Fatal("Failed to scan emergency values into variables: " + err.Error())
	}
	return emergencyAmount, emergencyMax
}

// used in different module so its grayed out
func IncreaseEmergencyFund(amount float32) {
	emergencyAmount, _ := GetEmergencyData()
	// make adding to this a transaction, lowering spending money
	AddTransaction(&Transaction{Amount: amount, Description: "Emergency Fund: Refill fund"})
	_, err := Database.Exec("UPDATE Variables SET emergencyAmount = ?;", emergencyAmount+amount)
	if err != nil {
		log.Fatal("Error updating emergency values into variables table:" + err.Error())
	}
}

// used in different module so its grayed out
// UpdateMaxEmergencyFund will be ran in finance.go every month before calculating
// how much to put inside this account to calculate 6 months worth of expenses, and have
// the emergency fund cover that.
func UpdateMaxEmergencyFund() {
	maxAmount := GetMonthlyExpenses() * 6
	// will change for all rows because no 'where' clause, but there is only a single row
	_, err := Database.Exec("UPDATE Variables SET emergencyMax = ?;", maxAmount)
	if err != nil {
		log.Fatal("Error updating max emergency amount in Variables table: " + err.Error())
	}
}

func SpendEmergencyFund(amount float32) (enough bool) {
	_, amountSaved := GetEmergencyData()
	if amount > amountSaved {
		return false
	}
	_, err := Database.Exec("UPDATE Variables SET emergencyAmount = ?;", amountSaved-amount)
	if err != nil {
		log.Fatal("Error updating emergencyAmount in variables table:" + err.Error())
	}
	return true
}
