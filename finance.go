package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var spendingMoney float32 = 0
var estimatedSpendingMoney float32 = 0
var estimatedIncome float32 = 0
var totalChangeThisMonth float32 = 0

func queueMonthlyTask() {
	for {
		monthlyTask()
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())

	}
}

func monthlyTask() {

}

// When the program first comes online, calculate the spending money based on the transactions
// This is to prevent any desyncs from not being online during the start of the month or other
func calculateSpendingMoney() {

}
