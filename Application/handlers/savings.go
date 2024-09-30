package handlers

import (
	"fmt"

	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleSavings() {
	for {
		fmt.Printf("Welcome to savings!\nYour amount to contribute this month is: %v\nWhat would you like to do?\n1)\tChange amount saved per month", db.GetSavingsPerMonth())
		response, exit := utils.GetUserResponse("")
		if exit {
			return
		}
		if response == "1" {
			fResponse, exit := utils.GetUserResponseFloat("What would you like to change your monthly savings contribution to?")
			if exit {
				return
			}
			if fResponse < 0.0 {
				fmt.Println("The amount you plan to contribute to savings every month should not be a negative number")
				continue
			}
			db.SetSavingsPerMonth(fResponse)
			return
		}
	}

}
