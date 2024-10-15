package main

import (
	"fmt"
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
	_, emergencyAmount := db.GetEmergencyData()

	for {
		// TODO: When you edit or delete a transaction, make it so it updates everything properly
		response, exit := utils.GetUserResponse(`Welcome to Finance!
		Spending money is: $%v
		Your emergency fund should be at: $%v
		You should add $%v to your savings account for last month
		What would you like to do?
				1) Add a transaction
				2) Display and edit all transactions
				3) View and edit monthly expenses
				4) View and edit goals
				5) Manage your emergency fund
				6) Manage your savings
				7) Change expected monthly income
				8) Pass a month by for testing`, db.GetEstimatedSpendingMoney(), emergencyAmount, db.GetAmountToSaveThisMonth())

		if exit {
			return
		}
		switch response {
		case "1":
			handlers.HandleAddTransaction()
		case "2":
			handlers.HandleDisplayEditTransactions()
		case "3":
			handlers.HandleViewAndEditMonthlyExpenses()
		case "4":
			handlers.HandleViewAndEditGoal()
		case "5":
			handlers.HandleEmergencyFund()
		case "6":
			handlers.HandleSavings()
		case "7":
			handlers.HandleChangeExpectedIncome()
		case "8":
			//pass a month for testing
			shared.MonthlyTask()

		default:
			fmt.Println("Invalid input")
			continue
		}
	}

}

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() float32 {
	transactions := db.GetAllTransactions(nil, nil)

	var spendingMoney float32 = 0.0
	for i := 0; i < len(transactions); i++ {
		spendingMoney += transactions[i].Amount
	}
	return spendingMoney
}
