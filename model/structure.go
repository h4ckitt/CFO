package model

import "time"

type Spending struct {
	MessageID int
	Amount    int
	Category  string
	Note      string
	CreatedAt time.Time
}

type CategorySpending struct {
	Category      string
	TotalSpending float64
}
