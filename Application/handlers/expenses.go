package handlers

import (
	"fmt"
	"math"

	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

var selectedExpense db.MonthlyExpense
var selectedGoal db.Goal

func HandleViewAndEditMonthlyExpenses() {
	for {
		var monthlyExpenses []db.MonthlyExpense = db.GetAllMonthlyExpensesStructs()

		// optionStrings := []string{}
		// for name, amount := range monthlyExpenses {
		// 	optionStrings = append(optionStrings, fmt.Sprintf("Name: %v\tAmount: %v", name, amount))
		// }
		var selectedExpensePtr *db.MonthlyExpense
		var exit bool
		for {
			selectedExpensePtr, exit = utils.SelectRecordOrCreate(monthlyExpenses, handleAddNewMonthlyExpense)
			if exit {
				return
			}
			break
		}
		// if the above method call wasnt exited or created new record
		if selectedExpensePtr == nil {
			return
		}
		selectedExpense = *selectedExpensePtr
		editMonthlyExpense()
	}

}

func editMonthlyExpense() {
	options := []string{"Change the name", "Change the amount", "Delete the expense"}
	methods := []func(){handleChangeExpenseName, handleChangeExpenseMonthlyAmount, handleDeleteExpense}
	var exit bool = false
	for !exit {
		exit = utils.PromptAndHandle("You have selected %v. Please select an option:", options, methods)
	}
}

func handleChangeExpenseName() {
	response, exit := utils.GetUserResponse("Please enter the new name for this expense: ")
	if exit {
		return
	}
	selectedExpense.UpdateExpenseName(response)
	fmt.Println("The expense name has been updated to " + response)
}

func handleChangeExpenseMonthlyAmount() {
	parsedFloat, exit := utils.GetUserResponseFloat("Please enter the new monthly amount for this expense: ", utils.Positive)
	if exit {
		return
	}
	// values should always be positive, they are assumed to be a negative transaction
	parsedFloat = float32(math.Abs(float64(parsedFloat)))

	oldAmount := selectedExpense.Amount
	newAmount := float32(parsedFloat)
	selectedExpense.UpdateExpenseAmount(newAmount)
	amountChanged := newAmount - oldAmount
	// updated the estimated spending money
	db.SetEstimatedSpendingMoney(db.GetEstimatedSpendingMoney() - amountChanged)
}

func handleDeleteExpense() {
	methods := []func(){
		func() {
			fmt.Printf("The %v monthly expense has been deleted!\n", selectedExpense.Name)
		},
		func() {
			db.SetEstimatedSpendingMoney(db.GetEstimatedSpendingMoney() + selectedExpense.Amount)
			fmt.Printf("The %v monthly expense has been deleted and $%v has been added to your estimated spending money!\n",
				selectedExpense.Name, selectedExpense.Amount)
		},
	}

	var exit bool = false
	for !exit {
		exit = utils.PromptAndHandle("Was the payment made already this month?", []string{"Yes", "No"}, methods)
	}

	selectedExpense.Delete()
}

func handleAddNewMonthlyExpense() {
	expense := db.MonthlyExpense{}
	var exit bool
	expense.Name, exit = utils.GetUserResponse("Please enter the name for the new expense: ")
	if exit {
		return
	}
	parsedFloat, exit := utils.GetUserResponseFloat("Please enter the monthly amount for the new expense: ", utils.Positive)
	if exit {
		return
	}

	// values should always be positive, they are assumed to be a negative transaction
	parsedFloat = float32(math.Abs(float64(parsedFloat)))
	expense.Amount = parsedFloat

	db.AddMonthlyExpense(expense)
	// next month it will automatically calculate this, but for this month we just
	// adjust the estimated spending money
	db.SetEstimatedSpendingMoney(db.GetEstimatedSpendingMoney() - expense.Amount)
}

func HandleViewAndEditGoal() {
	for {

		fmt.Println("Your amount per month will be deducted from your estimated spending money, if the estimated spending money cannot support that amount, it will split the amount you do have among all goals")
		var selectedGoalPtr *db.Goal
		var exit bool
		for {
			selectedGoalPtr, exit = utils.SelectRecordOrCreate(db.GetAllGoalStructs(), handleAddNewGoal)
			if exit {
				return
			}
			if selectedGoalPtr == nil {
				return
			}
			selectedGoal = *selectedGoalPtr
			break
		}
		options := []string{"Edit goal values", "Contribute one time payment to goal", "Delete goal"}
		methods := []func(){editGoal, contributeToGoal, selectedGoal.DeleteGoal}
		utils.PromptAndHandle("What would you like to do with your %v goal?", options, methods, selectedGoal.Name)
	}
}

// have two types of goals. Goals you set a date and it tells you how much money to put towards it per month
// and goals you set a monthly amount you can afford and it will tell you the date you will have enough
func handleAddNewGoal() {
	goal := db.Goal{}
	var exit bool
	goal.Name, exit = utils.GetUserResponse(`What would you like to name this goal?`)
	if exit {
		return
	}
	goal.Amount, exit = utils.GetUserResponseFloat("How much must be saved to complete this goal?", utils.Positive)
	if exit {
		return
	}

	for {
		response, exit := utils.GetUserResponse(`How would you like to create this goal?
							1) Set an amount you can afford per month
							2) Set an amount of months you would like the goal met by`)
		if exit {
			return
		}
		switch response {
		case "1":
			for {
				goal.AmountPerMonth, exit = utils.GetUserResponseFloat("What is the amount you can afford per month?", utils.Positive)
				if exit {
					return
				}
				goal.PopulateMonthsLeft()
				break
			}
		case "2":
			goal.MonthsLeft, exit = utils.GetUserResponseInt("In how many months would you like the goal complete?", utils.Positive)
			if exit {
				return
			}
			goal.PopulateAmountPerMonth()
		default:
			fmt.Println("Invalid Input")
		}
		break
	}
	db.AddGoal(&goal)
	fmt.Printf("Your goal was successfully created, you will save $%v per month for %v months\n", goal.AmountPerMonth, goal.MonthsLeft)

}

//func getDateFromUser() (date time.Time, exit bool) {
//	returnTime := time.Time{}
//	for {
//		yearInt, exit := utils.GetUserResponseInt("What year would you like the goal to be met by?")
//		if exit {
//			return returnTime, true
//		}
//		monthInt, exit := utils.GetUserResponseInt("What month would you like the goal to be met by?")
//		if exit {
//			return returnTime, true
//		}
//		dayInt, exit := utils.GetUserResponseInt("What day would you like the goal to be met by?")
//		if exit {
//			return returnTime, true
//		}
//		returnTime = time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
//		return returnTime, false
//	}
//}

func editGoal() {
	response, exit := utils.GetUserResponse(`Would you like to edit:
												1) Goal name
												2) Goal amount
												3) Months to goal completion
												4) Amount per month towards the goal`)
	if exit {
		return
	}
	switch response {
	case "1":
		response, exit := utils.GetUserResponse("What would you like the goals new name to be?")
		if exit {
			return
		}
		selectedGoal.UpdateGoalName(response)
	case "2":
		var changeMonthlyPayments bool
		response, exit := utils.GetUserResponseFloat("What would you like the new goal amount to be?", utils.Positive)
		if exit {
			return
		}
		for {
			methodResponse, exit := utils.GetUserResponse("Do you want to:\n1) Automatically change the date of goal completion\n2) Automatically change the monthly payment")
			if exit {
				return
			}
			if methodResponse == "1" {
				changeMonthlyPayments = false
			} else if methodResponse == "2" {
				changeMonthlyPayments = true
			} else {
				fmt.Println("Invalid Input")
				continue
			}
			break
		}
		selectedGoal.UpdateGoalAmount(response, changeMonthlyPayments)
		fmt.Printf("Successfully updated goal.\nMontly payments: $%v\nMonthly payments left: %v", selectedGoal.AmountPerMonth, selectedGoal.MonthsLeft)
	case "3":
		months, exit := utils.GetUserResponseInt("How many months from now would you like the goal complete instead?", utils.Positive)
		if exit {
			return
		}
		selectedGoal.UpdateMonthsLeft(months)
		fmt.Printf("Your monthly contribution will now be %v to achieve your goal in %v months", selectedGoal.AmountPerMonth, selectedGoal.MonthsLeft)
	case "4":
		response, exit := utils.GetUserResponseFloat("What would you like the new monthly payment to be?", utils.Positive)
		if exit {
			return
		}
		selectedGoal.UpdateMonthly(response)
		fmt.Printf("Your goal is now set to be completed in %v months\n", selectedGoal.MonthsLeft)

	}
}

func contributeToGoal() {
	response, exit := utils.GetUserResponseFloat("How much would you like to contribute to the goal? (This will be deducted from your spending money)", utils.Positive)
	if exit {
		return
	}
	selectedGoal.Contribute(response)
}
