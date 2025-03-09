package handlers

import (
	"fmt"
	"strconv"

	"github.com/ticotrem/finance/shared/db"

	"github.com/ticotrem/finance/shared/utils"
)

func HandleAddTransaction() {
	for {
		amountString, exit := utils.GetUserResponse("What is the amount of the transaction?")
		if exit {
			return
		}
		parsedFloat, err := strconv.ParseFloat(amountString, 32)
		amount := float32(parsedFloat)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		db.AddTransaction(&db.Transaction{Amount: amount, Date: utils.CurrentTime().UTC(), Description: "(Custom) User Added"})
		fmt.Println("Your transaction has successfully been added to the database!")
		// if the transaction is positive, we just need to add it to the database as above.
		// if it is negative, we will lower the estimated spending money accordingly
		if amount < 0 {
			db.SetEstimatedSpendingMoney(db.GetEstimatedSpendingMoney() + amount) // adding a negative is still negative
		}
		break
	}
}

var selectedTransaction db.Transaction

func HandleDisplayEditTransactions() {
	for {
		now := utils.CurrentTime().AddDate(0, 0, 1).UTC()
		dBegin := now.AddDate(0, -1, -2)
		// get all transactions within the past month
		transactions := db.GetAllTransactions(&dBegin, &now)

		var parsedInt int
		//for i := 0; i < len(transactions); i++ {
		//	fmt.Printf("%v:\tAmount: %v\t Date: %v\tDescription: %v\n", i+1, utils.GetMoneyString(transactions[i].Amount), transactions[i].Date.Local().Format(time.DateTime), transactions[i].Description)
		//}
		selectedTransactionPtr, exit := utils.SelectRecordOrCreate(transactions, HandleAddTransaction)
		if exit {
			return
		}

		if selectedTransactionPtr == nil {
			return
		}
		selectedTransaction = *selectedTransactionPtr

		options := []string{"Edit the transaction value", "Delete the transaction"}
		functions := []func(){handleEditTransaction, selectedTransaction.Delete}
		utils.PromptAndHandle("Transaction %v:\\tAmount: %v\\t Date: %v\\t was selected.\\nWould you like to:", options, functions, parsedInt, selectedTransaction.Amount, selectedTransaction.Date)
	}

}

func handleEditTransaction() {
	for {
		response, exit := utils.GetUserResponse("Please enter the new amount for this transaction: ")
		if exit {
			return
		}
		parsedFloat, err := strconv.ParseFloat(response, 32)
		if err != nil {
			fmt.Println("The value could not be converted in to a float!")
			continue
		}
		selectedTransaction.Amount = float32(parsedFloat)
		db.UpdateTransaction(&selectedTransaction)
		break
	}
}
