package handlers

import (
	"fmt"
	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

var maxAmount float32
var monthlyPaymentWhenUnder float32

func HandleEmergencyFund() {

	// emergency fund should cover 6 months of your expenses (also can be used for random surprise payments

	emergencyAmount, emergencyMax := db.GetEmergencyData()
	fmt.Printf("Welcome to to your emergency fund.\nCurrent fund is $%v/$%v\n"+
		"If your emergency fund is not full, half your new monthly spending money will "+
		"go towards it until filled.\n", emergencyMax, emergencyAmount)

	response, exit := utils.GetUserResponse(`What would you like to do?
			1) Spend Emergency Fund`)
	if exit {
		return
	}
	switch response {
	case "1":
		// It is not a transaction record to spend the emergency fund,
		// it is a transaction to fill it
		response, exit := utils.GetUserResponseFloat(`How much would you like to spend? (Note this will not change spending money)`)
		if exit {
			return
		}
		enough := db.SpendEmergencyFund(response)
		if !enough {
			fmt.Println("There is not enough money in your emergency fund. I hope you have savings!")
		}
	}

}
