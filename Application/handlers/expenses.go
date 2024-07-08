package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleViewAndEditMonthlyExpenses() {
	monthlyExpenses := shared.GetAllMonthlyExpensesStructs()
	for i := 0; i < len(monthlyExpenses); i++ {
		fmt.Printf("Monthly expense %v:\nName: %v\nAmount: %v\n", i+1, monthlyExpenses[i].Name, monthlyExpenses[i].Amount)
	}
	response, exit := utils.GetUserResponse("Enter the number of the expense you would like to edit, or 'C' to create a new one")
	if exit {
		return
	}
	if strings.ToLower(response) == "c" {
		handleAddNewMonthlyExpense()
	} else {
		editMonthlyExpense(response, monthlyExpenses)
	}

}

func editMonthlyExpense(response string, monthlyExpenses []shared.MonthlyExpense) {
	var selectedExpense shared.MonthlyExpense
	for {
		parsedInt, err := strconv.Atoi(response)
		if err != nil || parsedInt < 0 || parsedInt > len(monthlyExpenses) {
			fmt.Println("Invalid input")
			continue
		}

		parsedInt, err = strconv.Atoi(response)
		if err != nil {
			fmt.Println("Invalid Input")
			continue
		}
		selectedExpense = monthlyExpenses[parsedInt-1]
		break
	}

	fmt.Printf(`You have selected %v. Please select an option:
	1) Change the name
	2) Change the amount
	3) Delete the expense`, selectedExpense.Name)

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
			_, err := shared.Database.Exec("UPDATE MonthlyExpense SET name = ? WHERE id = ?", response, selectedExpense.Id)
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
			oldAmount := selectedExpense.Amount
			newAmount := float32(float64bit)
			_, err = shared.Database.Exec("UPDATE MonthlyExpense SET amount = ? WHERE id = ?", newAmount, selectedExpense.Id)
			if err != nil {
				log.Fatal("Failed update database expense amount: " + err.Error())
			}
			amountChanged := newAmount - oldAmount
			// updated the estimated spending money
			shared.SetEstimatedSpendingMoney(shared.GetEstimatedSpendingMoney() + amountChanged)
		case "3":

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
