package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleViewAndEditMonthlyExpenses() {
	monthlyExpenses := shared.GetAllMonthlyExpensesStructs()
	for i := 0; i < len(monthlyExpenses); i++ {
		fmt.Printf("%v:\tName: %v\tAmount: %v\n", i+1, monthlyExpenses[i].Name, monthlyExpenses[i].Amount)
	}
	response, createNew, exit := utils.CreateNewOrInt("Enter the number of the expense you would like to edit, or 'c' to create a new one", 1, len(monthlyExpenses))
	if exit {
		return
	} else if createNew {
		handleAddNewMonthlyExpense()
	} else {
		editMonthlyExpense(monthlyExpenses[response-1])
	}

}

func editMonthlyExpense(monthlyExpense shared.MonthlyExpense) {

	fmt.Printf(`You have selected %v. Please select an option:
	1) Change the name
	2) Change the amount
	3) Delete the expense`, monthlyExpense.Name)

	for {
		response, exit := utils.GetUserResponse("")
		if exit {
			return
		}
		switch response {
		case "1":
			response, exit := utils.GetUserResponse("Please enter the new name for this expense: ")
			if exit {
				return
			}
			_, err := shared.Database.Exec("UPDATE MonthlyExpenses SET name = ? WHERE id = ?;", response, monthlyExpense.Id)
			if err != nil {
				log.Fatal("Failed to update the expense name: " + err.Error())
			}
			fmt.Println("The expense name has been updated to " + response)
		case "2":
			newMonthlyResponse, exit := utils.GetUserResponse("Please enter the new monthly amount for this expense: ")
			if exit {
				return
			}
			float64bit, err := strconv.ParseFloat(newMonthlyResponse, 32)
			if err != nil {
				fmt.Println("The value could not be converted in to a float!")
				continue
			}
			oldAmount := monthlyExpense.Amount
			newAmount := float32(float64bit)
			_, err = shared.Database.Exec("UPDATE MonthlyExpenses SET amount = ? WHERE id = ?", newAmount, monthlyExpense.Id)
			if err != nil {
				log.Fatal("Failed update database expense amount: " + err.Error())
			}
			amountChanged := newAmount - oldAmount
			// updated the estimated spending money
			shared.SetEstimatedSpendingMoney(shared.GetEstimatedSpendingMoney() + amountChanged)
		case "3":
			//TODO:
		default:
			fmt.Println("Invalid input")
			continue
		}
		break
	}

}

func handleAddNewMonthlyExpense() {
	expense := shared.MonthlyExpense{}
	for {
		var exit bool
		expense.Name, exit = utils.GetUserResponse("Please enter the name for the new expense: ")
		if exit {
			return
		}
		amountString, exit := utils.GetUserResponse("Please enter the monthly amount for the new expense: ")
		if exit {
			return
		}
		parsedFloat, err := strconv.ParseFloat(amountString, 32)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		expense.Amount = float32(parsedFloat)
		break
	}

	shared.AddMonthlyExpense(expense)
	// next month it will automatically calculate this, but for this month we just
	// adjust the estimated spending money
	shared.SetEstimatedSpendingMoney(shared.GetEstimatedSpendingMoney() + expense.Amount)

}

func HandleViewAndEditGoal() {
	monthlyGoals := shared.GetAllGoalStructs()
	for i := 0; i < len(monthlyGoals); i++ {
		fmt.Printf("%v:\tName: %v\tAmount: %v/%v\tAmount Per Month: %v\n", i+1, monthlyGoals[i].Name, monthlyGoals[i].AmountSaved, monthlyGoals[i].Amount, monthlyGoals[i].AmountPerMonth)
	}
	response, createNew, exit := utils.CreateNewOrInt("Enter the number of the goal you would like to edit, or 'c' to create a new one", 1, len(monthlyGoals))
	if exit {
		return
	} else if createNew {
		handleAddNewGoal()
	} else {
		manageGoal(monthlyGoals[response-1])
	}

}

// have two types of goals. Goals you set a date and it tells you how much money to put towards it per month
// and goals you set a monthly amount you can afford and it will tell you the date you will have enough
func handleAddNewGoal() {
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
				shared.
			}
		default:
			fmt.Println("Invalid Input")
		}
		break
	}
	shared.AddGoal(&goal)

}

func manageGoal(goal shared.Goal) {
	response, exit := utils.GetUserResponse(`What would you like to do with your %v goal?
										1) Edit goal values
										2) Contribute one time payment to goal
										3) Delete goal`, goal.Name)
	if exit {
		return
	}
	switch response {
	case "1":
		editGoal(goal)
	case "2":
		contributeToGoal(goal)
	case "3":
		goal.DeleteGoal()
	}
}


func editGoal(goal shared.Goal) {
	response, exit := utils.GetUserResponse(`Would you like to edit:
												1) Goal name
												2) Goal completion date
												3) Amount per month towards the goal)
}
