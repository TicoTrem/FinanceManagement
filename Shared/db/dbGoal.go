package db

import (
	"fmt"
	"log"
	"math"
	"time"
)

type Goal struct {
	Id             int
	Name           string
	Amount         float32
	AmountSaved    float32
	AmountPerMonth float32
	MonthsLeft     int
}

func AddGoal(goal *Goal) {
	_, err := Database.Exec("INSERT INTO Goals (name, amount, amountSaved, monthsLeft) VALUES (?, ?, ?, ?);", goal.Name, goal.Amount, goal.AmountSaved, goal.MonthsLeft)
	if err != nil {
		log.Fatal("Error inserting goal in to the database:" + err.Error())
	}
}

func GetAllGoalStructs() []Goal {
	var goals []Goal
	rows, err := Database.Query("SELECT * FROM Goals;")
	if err != nil {
		log.Fatal("Querying all goals failed" + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var goal Goal
		err = rows.Scan(&goal.Id, &goal.Name, &goal.Amount, &goal.AmountSaved, &goal.MonthsLeft)
		if err != nil {
			log.Fatal("Failed to scan goal into goal struct:" + err.Error())
		}
		// calculate and assign the amount per month attribute
		goal.PopulateAmountPerMonth()

		goals = append(goals, goal)
	}
	return goals
}

// populates the AmountPerMonth attribute of the struct, assuming it has Amount and DateComplete specified
func (goal *Goal) PopulateAmountPerMonth() {
	if goal.Amount == 0 || goal.MonthsLeft == 0 {
		log.Fatal("You cannot populate the amount per month without having the Amount and MonthsLeft fields")
	}

	if goal.MonthsLeft <= 0 {
		goal.AmountPerMonth = float32(0)
	} else {
		goal.AmountPerMonth = (goal.Amount - goal.AmountSaved) / float32(goal.MonthsLeft)
	}
}

// populates the DateComplete attribute of the struct, assuming it has AmountPerMonth and Amount specified
func (goal *Goal) PopulateMonthsLeft() {
	if goal.Amount == 0 || goal.AmountPerMonth == 0 {
		log.Fatal("You cannot populate the months left without having the Amount and MonthsLeft attributes")
	}
	// round up, so if it takes 3.1 months, we will give it 4 months (as they can't afford it in the previous month
	var months int = int(math.Ceil(float64(goal.Amount / goal.AmountPerMonth)))
	goal.MonthsLeft = months
}

// overrideValue should be a positive float32, it will override the monthly if it is not 0
func (goal *Goal) SaveMonthlyAmount(overrideValue float32) {
	var amountToTransact float32 = -goal.AmountPerMonth
	if overrideValue > 0.0 {
		amountToTransact = -overrideValue
	}
	AddTransaction(&Transaction{Amount: amountToTransact, Date: time.Now().AddDate(0, 0, -1), Description: fmt.Sprintf("Goal: %v monthly savings", goal.Name)})
	_, err := Database.Exec("UPDATE Goals SET amountSaved = ? WHERE id = ?;", goal.AmountSaved+amountToTransact, goal.Id)
	if err != nil {
		log.Fatal("Error updating goal in database: " + err.Error())
	}
}

func (goal *Goal) UpdateGoalName(name string) {
	_, err := Database.Exec("UPDATE Goals SET name = ? WHERE id = ?;", name, goal.Id)
	if err != nil {
		log.Fatal("Failed to update the goal name: " + err.Error())
	}
}

// UpdateGoalAmount updates the goal amount in the database, and either increases monthly payment to compensate or increases
// the time to complete the goal as specified from the payMoreMonthly boolean
func (goal *Goal) UpdateGoalAmount(amount float32, payMoreMonthly bool) {
	goal.Amount = amount
	var err error
	// adjust the monthly amount
	if payMoreMonthly {
		// this will calculate and apply a new value since goal.Amount changed
		goal.PopulateAmountPerMonth()
		// we do not need to update the AmountPerMonth in the database as its only something calculated
		_, err = Database.Exec("UPDATE MonthlyExpenses SET amount = ? WHERE id = ?;", amount, goal.Id)
	} else { // adjust the time to completion
		// this will calculate and apply a new value since goal.Amount changed
		goal.PopulateMonthsLeft()
		_, err = Database.Exec("UPDATE MonthlyExpenses SET amount = ? AND dateComplete = ? WHERE id = ?;", amount, goal.DateComplete, goal.Id)
	}
	if err != nil {
		log.Fatal("Failed to update the expense amount: " + err.Error())
	}
}
func (goal *Goal) UpdateGoalDate(date time.Time) {
	_, err := Database.Exec("UPDATE Goals SET dateComplete = ? WHERE id = ?;", date, goal.Id)
	if err != nil {
		log.Fatal("Failed to update the goal date in database: " + err.Error())
	}
	// update the amount per month with new completion date
	goal.PopulateAmountPerMonth()
}
func (goal *Goal) UpdateGoalMonthly(amountPerMonth float32) {
	goal.AmountPerMonth = amountPerMonth
	// updates the date it will be completed by based on the amount per month and the unchanged amount
	goal.PopulateMonthsLeft()
	// update date in database
	_, err := Database.Exec("UPDATE Goals SET dateComplete = ? WHERE id = ?;", goal.DateComplete, goal.Id)
	if err != nil {
		log.Fatal("Failed to update the dateComplete based on the new amountPerMonth in database: " + err.Error())
	}
}
func (goal *Goal) DeleteGoal() {
	_, err := Database.Exec("DELETE FROM Goals WHERE id = ?;", goal.Id)
	if err != nil {
		log.Fatal("Error deleting goal from the database: " + err.Error())
	}
	*goal = Goal{} // make it zero value since it is now deleted from database
}
func (goal *Goal) Contribute(amountToContribute float32) {
	_, err := Database.Exec("UPDATE Goals SET amountSaved = ? WHERE id = ?;", goal.AmountSaved+amountToContribute, goal.Id)
	if err != nil {
		log.Fatal("Error contributing to goal in database: " + err.Error())
	}
}
