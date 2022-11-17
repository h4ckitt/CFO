package repository

import (
	"cfo/model"
)

type Repo interface {
	WaitForDB(addr, port string) bool
	SaveEntry(userId int, entry model.Spending) error
	RetrieveSpending(userId int, start, end string) ([]model.Spending, error)
	RetrieveSpendingByCategory(userId int, start, end string) ([]model.CategorySpending, error)
}
