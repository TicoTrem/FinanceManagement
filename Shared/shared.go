package shared

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/ticotrem/finance/shared/db"
	"github.com/ticotrem/finance/shared/utils"
)

var isTesting bool = true

// This function will setup the database and create the tables if they don't exist
func SetupDatabase() {

	dbase, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal(err)
	}

	var databaseName string
	// if isTesting {
	// 	databaseName = "FinanceTesting"
	// } else {
	// 	databaseName = "Finance"
	// }
	databaseName = "Finance"
	_, err = dbase.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", databaseName))

	if err != nil {
		log.Fatal(err)
	}
	dbase.Close()

	// Create the database object for real
	db.Database, err = sql.Open("mysql", "root:root@/Finance")
	if err != nil {
		log.Fatal(err)
	}

	createTables()

	// get the starting spending money (intensive operation)

	if db.Database.Ping() != nil {
		log.Fatal("Failed to ping database")
	}
}

// creates the tables needed for the application if they are not created already
// also populates the Variables table with a row containing all 0.0 if there is not already a row
func createTables() {
	_, err := db.Database.Exec(`CREATE TABLE IF NOT EXISTS Transactions (
		id INT AUTO_INCREMENT,
		amount FLOAT(16,2) NOT NULL,
		date DATETIME NOT NULL,
    	description VARCHAR(100) NOT NULL,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Database.Exec(`CREATE TABLE IF NOT EXISTS MonthlyExpenses (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		amount FLOAT(16,2) NOT NULL,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	// consider changing Goals table to store months left instead of dateComplete
	// TODO: Change goals to just keep track of months instead of a date
	_, err = db.Database.Exec(`CREATE TABLE IF NOT EXISTS Goals (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		amount FLOAT(16,2) NOT NULL,
    	amountSaved float(16,2) NOT NULL,
    	monthsLeft int NOT NULL DEFAULT 0,
		PRIMARY KEY(id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Database.Exec(`CREATE TABLE IF NOT EXISTS Variables (
	spendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedSpendingMoney FLOAT(16,2) DEFAULT 0.0,
	estimatedIncome FLOAT(16,2) DEFAULT 0.0,
    emergencyMax float(16,2) DEFAULT 0.0,
    emergencyAmount float(16,2) DEFAULT 0.0,
	emergencyFillFactor float(16, 2) DEFAULT 0.5,
	savingsPerMonth float(16,2) DEFAULT 0.0,
	amountToSaveThisMonth float(16,2) DEFAULT 0.0);`)
	if err != nil {
		log.Fatal(err)
	}
	row := db.Database.QueryRow(`SELECT COUNT(*) FROM Variables`)
	var count int
	err = row.Scan(&count)

	if err != nil {
		log.Fatal(err)
	}
	// there are no rows in the table
	if count == 0 {
		db.Database.Exec(`INSERT INTO Variables () VALUES ()`)
	}

}

func MonthlyTask() {

	// add the monthly transactions for that month and set the date as the last day of the previous month
	// this is so you can make changes to the prices of the monthly expenses and have it reflected in that months transactions
	// TODO: make sure the estimated spending money is updated when the user alters a monthly expense value, it is assumed to be active THAT MONTH
	// so it must be updated.

	// the above 2 are in the transactions, so we can now calculate how much
	// our money went up or down this month in total (no estimated values)

	emergencyAmount, emergencyMax := db.GetEmergencyData()
	addMonthlyTransactions(emergencyAmount, emergencyMax)

	netTransactionChange := calculateNetTransactionChange()

	// contribute to emergency fund
	if emergencyAmount < emergencyMax && netTransactionChange > 0 {
		// creates a transaction, lowering the estimated spending money by adding a transaction
		amountToIncrease := netTransactionChange * db.GetEmergencyFillFactor()
		difference := emergencyMax - emergencyAmount
		if amountToIncrease > difference {
			// only increase the difference
			amountToIncrease = difference
		}
		db.IncreaseEmergencyFund(amountToIncrease)
	}

	spendingMoney := db.GetSpendingMoney() + netTransactionChange

	// update it for last month, we later updated EstimatedSpendingMoney for THIS month.
	db.SetSpendingMoney(spendingMoney)

	// Set the estimated spending money value to the spending money, with next months predicted outcome
	// and deducting the set in stone monthly expenses. The expenses should be automatically registered as transactions because
	// otherwise you would not be able to lower the spending money when you make purchases
	db.SetEstimatedSpendingMoney(spendingMoney + db.GetExpectedMonthlyIncome() + db.GetMonthlyExpenses())

}

func addMonthlyTransactions(emergencyAmount float32, emergencyMax float32) {
	// add expenses as transactions
	var expenses []db.MonthlyExpense = db.GetAllMonthlyExpensesStructs()
	for i := 0; i < len(expenses); i++ {
		db.AddTransaction(&db.Transaction{Amount: -expenses[i].Amount, Date: utils.CurrentTime().AddDate(0, 0, -1), Description: fmt.Sprintf("(Expenses) $%v %v monthly payment", expenses[i].Amount, expenses[i].Name)})
	}

	// The emergency fund takes half of the netTransaction change if it is positive.
	var amountToAddToEmergency float32 = 0.0
	amountToSave := db.GetSavingsPerMonth()
	estimatedSpendingMoney := db.GetEstimatedSpendingMoney()
	if amountToSave >= estimatedSpendingMoney {
		// only assigning savings to this to calculate how much
		// to add to emergency fund later. This makes it so emergency fund will
		// not take more than we are estimated to have
		amountToSave = estimatedSpendingMoney
	}

	// this is adding transactions for the last day of last month, so tell people what they should add to their savings,
	// but that value is actually from last month
	amountToFillEmergencyFund := emergencyMax - emergencyAmount
	// if emergency is full
	if amountToFillEmergencyFund <= 0 {
		// add full savings amount to savings
		db.SetAmountToSaveThisMonth(float32(math.Max(float64(amountToSave), 0.0)))
		if amountToSave > 0 {
			db.AddTransaction(&db.Transaction{Amount: -amountToSave, Date: utils.CurrentTime().AddDate(0, 0, -1), Description: fmt.Sprintf("(Savings) $%v monthly contribution", amountToSave)})
		}
	} else { // else if the emergency is not full
		// if the savings amount fully covers filling the emergency fund
		if amountToSave > amountToFillEmergencyFund {
			// the difference (amount needed to fill fund) is added to emergencyFund
			amountToAddToEmergency += amountToFillEmergencyFund
			amountToSaveThisMonth := amountToSave - amountToFillEmergencyFund
			// used to show the user how much to add to savings account that month
			db.SetAmountToSaveThisMonth(amountToSaveThisMonth)
			if amountToSave > 0 {
				// the difference is removed from the amount added to savings
				db.AddTransaction(&db.Transaction{Amount: -amountToSave, Date: utils.CurrentTime().AddDate(0, 0, -1), Description: fmt.Sprintf("(Savings) $%v monthly contribution", amountToSaveThisMonth)})
				amountToSave = 0.0
			}
		} else { // add whatever you can to emergency without saving
			amountToAddToEmergency += amountToSave
			// nothing is added to savings
			db.SetAmountToSaveThisMonth(0)
		}
	}
	// fill emergency fund before savings
	db.IncreaseEmergencyFund(amountToAddToEmergency)

	// add goals as transactions
	// this happens after everything else because it is of least priority
	// consider implementing a priority system to make this simpler if I make another
	// project like this one, or ever rework this project
	goals := db.GetAllGoalStructs()
	var amountToSavePerGoal float32 = 0.0
	if estimatedSpendingMoney > 0 {
		var sum float32 = 0.0
		for i := 0; i < len(goals); i++ {
			sum += goals[i].AmountPerMonth
		}
		if estimatedSpendingMoney < sum {
			amountToSavePerGoal = sum / float32(len(goals))
		}
		for i := 0; i < len(goals); i++ {
			goal := goals[i]
			goal.SaveMonthlyAmount(amountToSavePerGoal)

		}
	}
}

// This will calculate the net transaction change (which includes income and expenses)
// Then this will change estimated spending money to this value
func calculateNetTransactionChange() float32 {

	// get the first of last month
	var firstOfLastMonth time.Time
	today := utils.CurrentTime()
	if !isTesting {
		if today.Day() != 1 {
			log.Fatal("The monthly task was ran on a day other than the 1st. Please fix this!")
		}
		firstOfLastMonth = today.AddDate(0, -1, 0)
	} else {
		// if testing, make first of last month actually equal yesterday
		firstOfLastMonth = today.AddDate(0, 0, -1)
	}

	// this will add a month, the subtract the amount of days, which takes us to the last day of the month
	lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -firstOfLastMonth.Day())

	var netTransactionChange float32 = 0.0

	lastMonthTransactions := db.GetAllTransactions(&firstOfLastMonth, &lastOfLastMonth)

	for i := 0; i < len(lastMonthTransactions); i++ {
		netTransactionChange += lastMonthTransactions[i].Amount
	}

	return netTransactionChange
}
