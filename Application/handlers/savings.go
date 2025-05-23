package handlers

import (
	"fmt"

	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleSavings() {
	fmt.Printf("Welcome to savings!\nYour amount to contribute this month is: %v\n", utils.GetMoneyString(db.GetSavingsPerMonth()))
	var exit bool = false
	for !exit {
		exit = utils.PromptAndHandle("What would you like to do?", []string{"Change amount saved per month", "Send some spending money to savings"}, []func(){handleChangeMonthlySavings, handleAddExtraSavings})
	}
}

func handleChangeMonthlySavings() {
	for {
		fResponse, exit := utils.GetUserResponseFloat("What would you like to change your monthly savings contribution to?", utils.Positive)
		if exit {
			return
		}
		valid := db.SetSavingsPerMonth(fResponse)
		if !valid {
			fmt.Println("You can't set your monthly savings contribution to be more than your current estimated spending money")
		}
		return
	}
}

func handleAddExtraSavings() {
	fResponse, exit := utils.GetUserResponseFloat("How much of your spending money would you like to add to savings?", utils.Positive)
	if exit {
		return
	}
	db.AddTransaction(&db.Transaction{Amount: -fResponse, Date: utils.CurrentTime().AddDate(0, 0, -1), Description: fmt.Sprintf("(Savings) $%v Additional Savings", fResponse)})
	fmt.Printf("Successfully added %v to savings", utils.GetMoneyString(fResponse))
}
