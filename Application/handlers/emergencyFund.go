package handlers

import (
	"fmt"

	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleEmergencyFund() {

	// emergency fund should cover 6 months of your expenses (also can be used for random surprise payments

	emergencyAmount, emergencyMax := db.GetEmergencyData()
	fmt.Printf("Welcome to to your emergency fund.\nCurrent fund is $%v/$%v\n"+
		"When not full, your emergency fund will intercept all funds going to savings, as well as %v of your monthly spending money\n"+
		"If your emergency fund is not full, half your new monthly spending money will "+
		"go towards it until filled.\n", emergencyAmount, emergencyMax, db.GetEmergencyFillFactor())

	utils.PromptAndHandle("What would you like to do?", []string{"Spend Emergency Fund", "Update Emergency Fund size", "Change what amount of spending money fills Emergency Fund instead"}, []func(){handleSpendEmergencyFund, handleUpdateMaxEmergencyFund, handleUpdateFillFactor})
}

func handleSpendEmergencyFund() {
	// It is not a transaction record to spend the emergency fund,
	// it is a transaction to fill it
	response, exit := utils.GetUserResponseFloat("How much would you like to spend? (Note: this will not change spending money)")
	if exit {
		return
	}
	enough := db.SpendEmergencyFund(response)
	if !enough {
		fmt.Println("There is not enough money in your emergency fund. I hope you have savings!")
	}
}

func handleUpdateMaxEmergencyFund() {
	response, exit := utils.GetUserResponseFloat("What's the new size you would like your Emergency Fund to be?")
	if exit {
		return
	}
	db.SetMaxEmergencyFund(response)
}

func handleUpdateFillFactor() {
	response, exit := utils.GetUserResponseFloat("How much of your spending money would you like to put towards filling your Emergency Fund when it is not full? (0.1 - 1.0)")
	if exit {
		return
	}
	if response > 1.0 || response < 0.1 {
		fmt.Println("You entered a value outside of the range 0.1 - 1.0")
		return
	}
	db.SetEmergencyFillFactor(response)

}
