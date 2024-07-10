package handlers

import (
	"fmt"
	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/utils"
	"strconv"
)

func HandleChangeExpectedIncome() {
	for {
		response, exit := utils.GetUserResponse("What is your new expected monthly income?")
		if exit {
			return
		}
		parsedFloat, err := strconv.ParseFloat(response, 32)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		income := float32(parsedFloat)
		oldMonthlyIncome := shared.GetExpectedMonthlyIncome()
		shared.SetEstimatedMonthlyIncome(income)
		fmt.Printf("Your expected monthly income has been set to %v. Estimations should be updated immediately!\n", income)
		// TODO: Find out why I solved this problem differently in the past, here I am just subtracting the old income and adding the
		// the new income to not have to recalculate anything
		shared.SetEstimatedSpendingMoney(shared.GetEstimatedSpendingMoney() - oldMonthlyIncome + income)
		break
	}
}
