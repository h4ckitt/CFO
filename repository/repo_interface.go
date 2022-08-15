package repository

import (
	"cfo/model"
)

type Repo interface {
	SaveEntry(userId int, entry model.Spending) error
	RetrieveSpending(userId int, start, end string) ([]model.Spending, error)
	RetrieveSpendingByCategory(userId int, category string) ([]model.CategorySpending, error)
}
