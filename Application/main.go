package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ticotrem/shared"
)

// TODO: Make sure we are closing all database connections (defer rows.close())
// I am havinig performance and battery issues from the MySQL container

func main() {

	db, err := sql.Open("mysql", "root:password@/Finance")
	shared.Database = db
	// You started the application before ever running the service

	// how can I make sure the service is currently running?
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(`Welcome to Finance!
				Spending money is: %v
				What would you like to do?
						1) Add a transaction
						2) Display and edit all transactions
						3) Change 'Expected' values
						4) View and edit monthly expenses
						5) Add a new goal to save up for`, shared.GetSpendingMoney())

	for {
		// TODO: When you edit or delete a transaction, make it so it updates everything properly
		var response string
		_, err = fmt.Scanln(&response)
		if err != nil {
			log.Fatal(err)
		}
		switch response {
		case "1":
			handleAddTransaction()
		case "2":
			printTransactions(shared.GetAllTransactions(nil, nil))
		case "3":
			handleChangeExpectedValues()
		case "4":
			handleViewAndEditMonthlyExpenses()
		case "5":
			handleAddNewGoal()
		default:
			fmt.Println("Invalid input")
			continue
		}
		break
	}

}

func handleViewAndEditMonthlyExpenses() {
	monthlyExpenses := shared.GetAllMonthlyExpensesStructs()
	for i := 0; i < len(monthlyExpenses); i++ {
		fmt.Printf("Monthly expense %v:\nName: %v\nAmount: %v\n", i+1, monthlyExpenses[i].Name, monthlyExpenses[i].Amount)
	}
	fmt.Println("Enter the number of the expense you would like to edit, or 'C' to create a new one")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	if strings.ToLower(response) == "c" {

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
		_, err := fmt.Scanln(&response)
		if err != nil {
			log.Fatal("Failed to read user input: " + err.Error())
		}
		switch response {
		case "1":
			fmt.Println("Please enter the new name for this expense: ")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				log.Fatal("Failed to read user input: " + err.Error())
			}
			_, err = shared.Database.Exec("UPDATE MonthlyExpense SET name = ? WHERE id = ?", response, selectedExpense.Id)
			if err != nil {
				log.Fatal("Failed to update the expense name: " + err.Error())
			}
			fmt.Println("The expense name has been updated to " + response)
		case "2":
			fmt.Println("Please enter the new monthly amount for this expense: ")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				log.Fatal("Failed to read user input: " + err.Error())
			}
			float64bit, err := strconv.ParseFloat(response, 32)
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

func handleAddTransaction() {
	for {
		fmt.Println("What is the amount of the transaction?")
		var amountString string
		_, err := fmt.Scanln(&amountString)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		parsedFloat, err := strconv.ParseFloat(amountString, 32)
		amount := float32(parsedFloat)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}
		addTransaction(amount, time.Now())
		break
	}

}

func addTransaction(amount float32, date time.Time) {
	transaction := shared.Transaction{Amount: amount, Date: date}
	query, err := shared.Database.Prepare("INSERT INTO Transactions (amount, date) VALUES (?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	result, err := query.Exec(transaction.Amount, transaction.Date)
	if err != nil {
		log.Fatal(err)
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There were %v rows inserted into the Transactions table\n", numRows)
}

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() float32 {
	transactions := shared.GetAllTransactions(nil, nil)

	var spendingMoney float32 = 0.0
	for i := 0; i < len(transactions); i++ {
		spendingMoney += float32(transactions[i].Amount)
	}
	return spendingMoney
}

func printTransactions(transactions []shared.Transaction) {
	for i := 0; i < len(transactions); i++ {
		fmt.Printf("Transaction %v:\nAmount: %v\n Date: %v\n", i+1, transactions[i].Amount, transactions[i].Date)
	}
}

func handleChangeExpectedValues() {
	fmt.Println("What is your expected monthly income?")
	var incomeString string
	_, err := fmt.Scanln(&incomeString)
	if err != nil {
		fmt.Println("Your change could not be made: " + err.Error())
		return
	}
	parsedIncome, err := strconv.ParseFloat(incomeString, 32)
	if err != nil {
		fmt.Println("Your change could not be made: " + err.Error())
		return
	}
	income := float32(parsedIncome)

	shared.SetExpectedMonthlyIncome(income)

	fmt.Printf("Your expected monthly income has been set to %v. Estimations should be updated immediately!", income)
}
