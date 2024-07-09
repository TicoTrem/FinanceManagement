package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ticotrem/finance/shared"
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

		shared.AddTransaction(&shared.Transaction{Amount: amount, Date: time.Now()})
		// if the transaction is positive, we just need to add it to the database as above.
		// if it is negative, we will lower the estimated spending money accordingly
		if amount < 0 {
			shared.SetEstimatedSpendingMoney(shared.GetEstimatedSpendingMoney() + amount) // adding a negative is still negative
		}
		break
	}

}

func HandleDisplayEditTransactions() {
	transactions := shared.GetAllTransactions(nil, nil)

	var selectedTransaction shared.Transaction
	var parsedInt int
	for i := 0; i < len(transactions); i++ {
		fmt.Printf("Transaction %v:\tAmount: %v\t Date: %v\n", i+1, transactions[i].Amount, transactions[i].Date)
	}
	for {
		response, exit := utils.GetUserResponse("If you would like to edit or delete one of these transactions, please enter the number of the transaction")
		if exit {
			return
		}
		parsedInt, err := strconv.Atoi(response)
		if err != nil || parsedInt < 0 || parsedInt > len(transactions) {
			fmt.Println("Invalid input")
			continue
		}
		selectedTransaction = transactions[parsedInt-1]
		break
	}
	for {
		response, exit := utils.GetUserResponse("Transaction %v:\tAmount: %v\t Date: %v\t was selected.\nWould you like to\n1) Edit the transaction value\n2) Delete the transaction\n",
			parsedInt, selectedTransaction.Amount, selectedTransaction.Date)
		if exit {
			return
		}
		switch response {
		case "1":
			handleEditTransaction(selectedTransaction)
		case "2":
			shared.DeleteTransaction(&selectedTransaction)
		default:
			fmt.Println("Invalid input")
			continue
		}
		break
	}

}

func handleEditTransaction(selectedTransaction shared.Transaction) {
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
		shared.UpdateTransaction(&selectedTransaction)
		break
	}
}
