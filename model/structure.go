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

type PendingReplyCache map[int]*Spending

func NewPendingReplyCache() PendingReplyCache {
	return make(map[int]*Spending)
}

func (p PendingReplyCache) CacheSpending(userID int, spend *Spending) {
	p[userID] = spend
}

func (p PendingReplyCache) RetrieveSpending(userID int) *Spending {
	return p[userID]
}

func (p PendingReplyCache) HasCachedSpending(userID int) bool {
	_, exists := p[userID]

	return exists
}

func (p PendingReplyCache) DeleteCachedSpending(userID int) {
	delete(p, userID)
}
