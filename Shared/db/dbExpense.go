package db

import (
	"fmt"
	"log"
)

type MonthlyExpense struct {
	Id     int
	Name   string
	Amount float32
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
		log.Fatal("Failed to update the expense amomlnghbunt: " + err.Error())
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
