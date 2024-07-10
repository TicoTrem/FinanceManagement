package shared

import "time"

type Transaction struct {
	Id          int
	Amount      float32
	Date        time.Time
	Description string
}

type MonthlyExpense struct {
	Id     int
	Name   string
	Amount float32
}

type Goal struct {
	Id             int
	Name           string
	Amount         float32
	AmountSaved    float32
	AmountPerMonth float32
	DateComplete   time.Time
}
