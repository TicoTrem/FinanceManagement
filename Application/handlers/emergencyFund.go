package handlers

import (
	"fmt"

	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleEmergencyFund() {

	// emergency fund should cover 6 months of your expenses (also can be used for random surprise payments

	emergencyAmount, emergencyMax := db.GetEmergencyData()
	if emergencyAmount <= 0 && emergencyMax <= 0 {
		handleSetEmergencyAmount()
		emergencyAmount, emergencyMax = db.GetEmergencyData()
	}
	fmt.Printf("Welcome to to your emergency fund.\nCurrent fund is %v/%v\n"+
		"When not full, your emergency fund will intercept all funds going to savings, as well as %.2f of your monthly spending money\n"+
		"If your emergency fund is not full, half your new monthly spending money will "+
		"go towards it until filled.\n", utils.GetMoneyString(emergencyAmount), utils.GetMoneyString(emergencyMax), db.GetEmergencyFillFactor())

	var exit bool = false
	for !exit {
		exit = utils.PromptAndHandle("What would you like to do?", []string{"Spend Emergency Fund", "Update Target Emergency Fund size", "Manually set amount in Emergency Fund", "Change what amount of spending money fills Emergency Fund instead"}, []func(){handleSpendEmergencyFund, handleUpdateMaxEmergencyFund, handleSetEmergencyAmount, handleUpdateFillFactor})
	}
}

func handleSpendEmergencyFund() {
	// It is not a transaction record to spend the emergency fund,
	// it is a transaction to fill it
	response, exit := utils.GetUserResponseFloat("How much would you like to spend? (Note: this will not change spending money)", utils.Positive)
	if exit {
		return
	}
	enough := db.SpendEmergencyFund(response)
	if !enough {
		fmt.Println("There is not enough money in your emergency fund. I hope you have savings!")
	}
}

// used to set up your emergency fund in the first place, shouldnt use this to deal with transactions or if interest brings your emergency fund higher
func handleSetEmergencyAmount() {
	response, exit := utils.GetUserResponseFloat("Set up your emergency fund! How much do you currently have in your emergency fund?", utils.Positive)
	if exit {
		return
	}
	db.SetEmergencyAmount(response)
	fmt.Printf("Your Emergency Fund balance has been set to: %v\n", utils.GetMoneyString(response))

}

func handleUpdateMaxEmergencyFund() {
	response, exit := utils.GetUserResponseFloat("What's the new size you would like your Emergency Fund to be?", utils.Positive)
	if exit {
		return
	}
	db.SetMaxEmergencyFund(response)
}

func handleUpdateFillFactor() {
	response, exit := utils.GetUserResponseFloat("How much of your spending money would you like to put towards filling your Emergency Fund when it is not full? (0.1 - 1.0)", utils.Positive)
	if exit {
		return
	}
	if response > 1.0 || response < 0.1 {
		fmt.Println("You entered a value outside of the range 0.1 - 1.0")
		return
	}
	db.SetEmergencyFillFactor(response)

}
