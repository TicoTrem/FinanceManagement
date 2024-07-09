package handlers

import (
	"fmt"
	"strconv"

	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/utils"
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
		shared.SetExpectedMonthlyIncome(income)
		fmt.Printf("Your expected monthly income has been set to %v. Estimations should be updated immediately!", income)
		break
	}
}

// have two types of goals. Goals you set a date and it tells you how much money to put towards it per month
// and goals you set a monthly amount you can afford and it will tell you the date you will have enough
func HandleAddNewGoal() {
	goal := shared.Goal{}
	goal.Name = utils.GetUserResponse(`What would you like to name this goal?`)

	for {
		goalAmountString, exit := utils.GetUserResponse("How much must be saved to complete this goal?")
		parsedFloat, err := strconv.ParseFloat(goalAmountString, 32)
		if err != nil {
			fmt.Println("Invalid Input")
			continue
		}
		goal.Amount = float32(parsedFloat)


	
		response, exit := utils.GetUserResponse(`How would you like to create this goal?
							1) Set an amount you can afford per month
							2) Set a date you would like the goal met by`)
		if exit {
			return
		}
		select response {
		case "1":
			response, exit = utils.GetUserResponse("What is the amount you can afford per month?")
			if exit {
				return
			}
			parsedFloat := strconv.ParseFloat(response, 32)
			goal.AmountPerMonth = float32(parsedFloat)
		}
			
	}
	
}
