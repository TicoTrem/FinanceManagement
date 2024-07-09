package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ticotrem/finance/handlers"
	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/utils"
)

// TODO: Make sure we are closing all database connections (defer rows.close())
// I am havinig performance and battery issues from the MySQL container

func main() {

	db, err := sql.Open("mysql", "root:password@/Finance")
	shared.Database = db
	// You started the application before ever running the service

	// how can I make sure the service is currently running?
	if err != nil {
		log.Fatal(err)
	}

	for {
		// TODO: When you edit or delete a transaction, make it so it updates everything properly
		response, exit := utils.GetUserResponse(`Welcome to Finance!
		Spending money is: %v
		What would you like to do?
				1) Add a transaction
				2) Display and edit all transactions
				3) Change 'Expected' values
				4) View and edit monthly expenses
				5) Add a new goal to save up for`, fmt.Sprint(shared.GetSpendingMoney()))
		if exit {
			return
		}
		switch response {
		case "1":
			handlers.HandleAddTransaction()
		case "2":
			handlers.HandleDisplayEditTransactions()
		case "3":
			handlers.HandleChangeExpectedIncome()
		case "4":
			handlers.HandleViewAndEditMonthlyExpenses()
		case "5":
			handlers.HandleAddNewGoal()
		default:
			fmt.Println("Invalid input")
			continue
		}
	}

}

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() float32 {
	transactions := shared.GetAllTransactions(nil, nil)

	var spendingMoney float32 = 0.0
	for i := 0; i < len(transactions); i++ {
		spendingMoney += transactions[i].Amount
	}
	return spendingMoney
}
