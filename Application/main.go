package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ticotrem/finance/handlers"
	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

// TODO: Make sure we are closing all database connections (defer rows.close())
// I am having performance and battery issues from the MySQL container

func main() {

	shared.SetupDatabase()

	// TODO: how can I make sure the service is currently running?

	options := []string{"Add a transaction", "Display and edit all transactions", "View and edit monthly expenses",
		"View and edit goals", "Manage your emergency fund", "Manage your savings",
		"Change expected monthly income", "Pass a month by for testing"}
	methods := []func(){handlers.HandleAddTransaction, handlers.HandleDisplayEditTransactions, handlers.HandleViewAndEditMonthlyExpenses,
		handlers.HandleViewAndEditGoal, handlers.HandleEmergencyFund, handlers.HandleSavings,
		handlers.HandleChangeExpectedIncome, shared.MonthlyTask}

	for {
		emergencyAmount, _ := db.GetEmergencyData()
		var exit bool = false
		for !exit {
			exit = utils.PromptAndHandle("Welcome to Finance!\nSpending money is: %v\nYour emergency fund should be at: $%v\n"+
				"You should add $%v to your savings account for last month\nWhat would you like to do?", options, methods,
				utils.GetMoneyString(db.GetEstimatedSpendingMoney()), emergencyAmount, db.GetAmountToSaveThisMonth())
		}

	}

}

// TODO When the service first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() float32 {
	transactions := db.GetAllTransactions(nil, nil)

	var spendingMoney float32 = 0.0
	for i := 0; i < len(transactions); i++ {
		spendingMoney += transactions[i].Amount
	}
	return spendingMoney
}
