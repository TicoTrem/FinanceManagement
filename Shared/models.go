package shared

import "time"

type Transaction struct {
	Id     int
	Amount float32
	Date   time.Time
}

type MonthlyExpense struct {
	Id     int
	Name   string
	Amount float32
}
