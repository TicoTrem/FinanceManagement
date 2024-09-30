package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ticotrem/finance/shared"
	"github.com/ticotrem/finance/shared/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ticotrem/finance/shared/utils"
)

func StartFinance() {

	shared.SetupDatabase()
	defer db.Database.Close()

	if db.GetEstimatedSpendingMoney() == 0 {
		db.SetEstimatedSpendingMoney(db.GetExpectedMonthlyIncome() - db.GetMonthlyExpenses())
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go queueMonthlyTask()
	fmt.Println("Service Running")
	waitGroup.Wait()

}

func queueMonthlyTask() {
	for {
		now := utils.CurrentTime()
		nextMonth := now.AddDate(0, 1, 0)
		firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		duration := firstOfNextMonth.Sub(now)

		timer := time.NewTimer(duration)

		// blocks execution until the timer expires (the start of the month)
		<-timer.C
		shared.MonthlyTask()
	}
}
