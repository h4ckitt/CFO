package service

import (
	"cfo/model"
)

type Service interface {
	SaveEntry(userID, messageID int, text string) error
	SendSpendingData(userID int, spending ...model.Spending) error
	SendGenericMessage(text string, userID int)
	RetrieveSpendingByDateRanges(userID int, ranges ...string) ([]model.Spending, error)
	RetrieveYesterdaySpending(userID int) ([]model.Spending, error)
	RetrieveThisWeekSpending(userID int) ([]model.Spending, error)
	RetrieveThisMonthSpending(userID int) ([]model.Spending, error)
}
