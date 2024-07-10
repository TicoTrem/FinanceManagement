package handlers

import (
	"fmt"
	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/utils"
	"strconv"
)

func HandleChangeExpectedIncome() {
	for {
		response, exit := utils.GetUserResponse("What is your new monthly expected income?")
		if exit {
			return
		}
		parsedFloat, err := strconv.ParseFloat(response, 32)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		income := float32(parsedFloat)
		shared.SetEstimatedMonthlyIncome(income)
		fmt.Printf("Your expected monthly income has been set to %v. Estimations should be updated immediately!\n", income)
		break
	}
}
