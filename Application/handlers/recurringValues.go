package handlers

import (
	"fmt"
	"strconv"
	"time"

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
	var exit bool
	goal.Name, exit = utils.GetUserResponse(`What would you like to name this goal?`)
	if exit {
		return
	}
	goal.Amount, exit = utils.GetUserResponseFloat("How much must be saved to complete this goal?")
	if exit {
		return
	}

	for {
		response, exit := utils.GetUserResponse(`How would you like to create this goal?
							1) Set an amount you can afford per month
							2) Set a date you would like the goal met by`)
		if exit {
			return
		}
		switch response {
		case "1":
			for {
				goal.AmountPerMonth, exit = utils.GetUserResponseFloat("What is the amount you can afford per month?")
				if exit {
					return
				}
			}
		case "2":
			for {
				yearInt, exit := utils.GetUserResponseInt("What year would you like the goal to be met by?")
				if exit {
					return
				}
				monthInt, exit := utils.GetUserResponseInt("What month would you like the goal to be met by?")
				if exit {
					return
				}
				dayInt, exit := utils.GetUserResponseInt("What day would you like the goal to be met by?")
				if exit {
					return
				}
				goal.DateComplete = time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
			}
		default:
			fmt.Println("Invalid Input")
		}
		break
	}
	shared.AddGoal(goal)

}
