package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ticotrem/shared"
)

func StartFinance() {

	// create database if it doesn't already exist

	go queueMonthlyTask()
}

func queueMonthlyTask() {
	for {
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		duration := firstOfNextMonth.Sub(now)

		timer := time.NewTimer(duration)

		// blocks execution until the timer expires (the start of the month)
		<-timer.C
		monthlyTask()
	}
}

func monthlyTask() {

	// add the monthly transactions for that month and set the date as the last day of the previous month
	// this is so you can make changes to the prices of the monthly expenses and have it reflected in that months transactions
	// TODO: make sure the estimated spending money is updated when the user alters a monthly expense value, it is assumed to be active THAT MONTH
	// so it must be updated.

	expenses := shared.GetAllMonthlyExpensesStructs()
	for i := 0; i < len(expenses); i++ {
		shared.AddTransaction(expenses[i].Amount, time.Now().AddDate(0, 0, -1))
	}

	netTransactionChange := calculateNetTransactionChange()

	spendingMoney := shared.GetSpendingMoney() + netTransactionChange
	shared.SetSpendingMoney(spendingMoney)

	// Set the estimated spending money value to the spending money, with next months predicted outcome
	// and deducting the set in stone monthly expenses. The expenses should be automatically registered as transactions because
	// otherwise you would not be able to lower the spending money when you make purchases
	shared.SetEstimatedSpendingMoney(spendingMoney + shared.GetExpectedMonthlyIncome() - shared.GetMonthlyExpenses())

}

// This will calculate the net transaction change (which includes income and expenses)
// Then this will change estimated spending money to this value
func calculateNetTransactionChange() float32 {

	// get the first of last month
	today := time.Now()
	if today.Day() != 1 {
		log.Fatal("The monthly task was ran on a day other than the 1st. Please fix this!")
	}

	firstOfLastMonth := today.AddDate(0, -1, 0)

	// this will add a month, the subtract the amount of days, which takes us to the last day of the month
	lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -firstOfLastMonth.Day())

	var netTransactionChange float32 = 0.0

	lastMonthTransactions := shared.GetAllTransactions(&firstOfLastMonth, &lastOfLastMonth)

	for i := 0; i < len(lastMonthTransactions); i++ {
		netTransactionChange += lastMonthTransactions[i].Amount
	}

	return netTransactionChange
}
