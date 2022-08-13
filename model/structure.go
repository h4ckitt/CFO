package model

import "time"

type Spending struct {
	MessageID int
	Amount    float32
	Category  string
	Note      string
	CreatedAt time.Time
}

type CategorySpending struct {
	Category      string
	TotalSpending float64
}
