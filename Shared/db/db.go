package db

import (
	"database/sql"
	"log"
	"math"

	"github.com/ticotrem/finance/shared/utils"
)

var Database *sql.DB

// TODO: store times in UTC time, then convert to show the user their local time

func GetExpectedMonthlyIncome() float32 {
	row := Database.QueryRow("SELECT estimatedIncome FROM Variables")
	var estimatedIncome float32
	row.Scan(&estimatedIncome)
	return estimatedIncome
}
func SetEstimatedMonthlyIncome(estimatedMonthlyIncome float32) {
	_, err := Database.Exec("UPDATE Variables SET estimatedIncome = ?", estimatedMonthlyIncome)
	if err != nil {
		log.Fatal("Failed to update the estimatedMonthlyIncome variable: " + err.Error())
	}
}

// This function will return the estimatedExpectedIncome variable from the Variables table
func GetEstimatedSpendingMoney() float32 {
	row := Database.QueryRow("SELECT estimatedSpendingMoney FROM Variables")
	var estimatedSpendingMoney float32
	err := row.Scan(&estimatedSpendingMoney)
	if err != nil {
		log.Fatal(err)
	}
	return estimatedSpendingMoney
}
func SetEstimatedSpendingMoney(estimatedSpendingMoney float32) {
	_, err := Database.Exec("UPDATE Variables SET estimatedSpendingMoney = ?", estimatedSpendingMoney)
	if err != nil {
		log.Fatal("Failed to update the estimatedSpendingMoney variable: " + err.Error())
	}
}

// This function will return the spendingMoney variable from the Variables table
func GetSpendingMoney() float32 {
	row := Database.QueryRow("SELECT spendingMoney FROM Variables")
	var spendingMoney float32
	err := row.Scan(&spendingMoney)
	if err != nil {
		log.Fatal(err)
	}
	return spendingMoney
}

func SetSpendingMoney(spendingMoney float32) {
	_, err := Database.Exec("UPDATE Variables SET spendingMoney = ?", spendingMoney)
	if err != nil {
		log.Fatal(err)
	}

}

func GetEmergencyData() (amount float32, max float32) {
	row := Database.QueryRow("SELECT emergencyMax, emergencyAmount FROM Variables")
	var emergencyMax float32
	var emergencyAmount float32
	err := row.Scan(&emergencyMax, &emergencyAmount)
	if err != nil {
		log.Fatal("Failed to scan emergency values into variables: " + err.Error())
	}
	return emergencyAmount, emergencyMax
}

func IncreaseEmergencyFund(amount float32) {
	if amount <= 0 {
		return
	}
	emergencyAmount, _ := GetEmergencyData()
	// make adding to this a transaction, lowering spending money
	AddTransaction(&Transaction{Amount: -amount, Date: utils.CurrentTime().AddDate(0, 0, -1), Description: "(Emergency Fund) Refill fund"})
	_, err := Database.Exec("UPDATE Variables SET emergencyAmount = ?;", emergencyAmount+amount)
	if err != nil {
		log.Fatal("Error updating emergency values into variables table:" + err.Error())
	}
}

// used in different module so its grayed out
// UpdateMaxEmergencyFund will be ran in finance.go every month before calculating
// how much to put inside this account to calculate 6 months worth of expenses, and have
// the emergency fund cover that.

func SetMaxEmergencyFund(maxAmount float32) {
	// Make sure the resulting number is positive
	maxAmount = float32(math.Abs(float64(maxAmount)))
	// will change for all rows because no 'where' clause, but there is only a single row
	_, err := Database.Exec("UPDATE Variables SET emergencyMax = ?;", maxAmount)
	if err != nil {
		log.Fatal("Error updating max emergency amount in Variables table: " + err.Error())
	}
}

func SetEmergencyFillFactor(fillFactor float32) {
	fillFactor = float32(math.Abs(float64(fillFactor)))
	_, err := Database.Exec("UPDATE Variables SET emergencyFillFactor = ?;", fillFactor)
	if err != nil {
		log.Fatal("Error updating emergencyFillFactor in Variables table: " + err.Error())
	}
}

func GetEmergencyFillFactor() float32 {
	row := Database.QueryRow("SELECT emergencyFillFactor FROM Variables")
	var emergencyFillFactor float32
	err := row.Scan(&emergencyFillFactor)
	if err != nil {
		log.Fatal("Failed to get emergencyFillFactor from Variables table: " + err.Error())
	}
	return emergencyFillFactor
}

func SpendEmergencyFund(amount float32) (enough bool) {
	_, amountSaved := GetEmergencyData()
	if amount > amountSaved {
		return false
	}
	_, err := Database.Exec("UPDATE Variables SET emergencyAmount = ?;", amountSaved-amount)
	if err != nil {
		log.Fatal("Error updating emergencyAmount in variables table:" + err.Error())
	}
	return true
}

func GetSavingsPerMonth() float32 {
	row := Database.QueryRow("SELECT savingsPerMonth FROM Variables")
	var savingsPerMonth float32
	err := row.Scan(&savingsPerMonth)
	if err != nil {
		log.Fatal("Failed to scan savingsPerMonth into variable: " + err.Error())
	}
	return savingsPerMonth
}

func SetSavingsPerMonth(savingsPerMonth float32) {

	_, err := Database.Exec("UPDATE Variables SET savingsPerMonth = ?;", savingsPerMonth)
	if err != nil {
		log.Fatal("Failed to update savingsPerMonth variable in database: " + err.Error())
	}
}

func GetAmountToSaveThisMonth() float32 {
	row := Database.QueryRow("SELECT amountToSaveThisMonth FROM Variables")
	var amountToSaveThisMonth float32
	err := row.Scan(&amountToSaveThisMonth)
	if err != nil {
		log.Fatal("Failed to scan amountToSaveThisMonth into variable: " + err.Error())
	}
	return amountToSaveThisMonth
}

func SetAmountToSaveThisMonth(amountToSaveThisMonth float32) {
	_, err := Database.Exec("UPDATE Variables SET amountToSaveThisMonth = ?;", amountToSaveThisMonth)
	if err != nil {
		log.Fatal("Failed to update amountToSaveThisMonth variable in database: " + err.Error())
	}
}

func GetMonthlyExpenses() float32 {
	var expenses []MonthlyExpense = GetAllMonthlyExpensesStructs()
	var differenceOfExpenses float32 = 0.0
	for i := 0; i < len(expenses); i++ {
		differenceOfExpenses -= expenses[i].Amount
	}
	return differenceOfExpenses
}
