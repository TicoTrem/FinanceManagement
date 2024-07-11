package handlers

import (
	"fmt"
	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

func HandleSavings() {
	fmt.Printf("Welcome to savings!\nYour current monthly contribution is: %v\n What would you like to do?\n1)Change amount saved per month", db.GetSavingsPerMonth())
	for {
		response, exit := utils.GetUserResponse("")
		if exit {
			return
		}
		if response == "1" {
			response, exit := utils.GetUserResponseFloat("What would you like to change your monthly savings contribution to?")
			if exit {
				return
			}
			db.SetSavingsPerMonth(response)
		}
	}

}
