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
	netTransactionChange := calculateNetTransactionChange()

	spendingMoney := shared.GetSpendingMoney() + netTransactionChange
	shared.SetSpendingMoney(spendingMoney)
	shared.SetEstimatedSpendingMoney(spendingMoney)
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
