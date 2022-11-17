package service

import (
	"cfo/model"
	"cfo/repository"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/h4ckitt/goTelegram"
)

var (
	decimalRegex = regexp.MustCompile(`^\d+(?:.\d{2})?$`)
	dateRegex    = regexp.MustCompile(`^\d{4}(?:-\d{2}){2}$`)
	dotRemover   = strings.NewReplacer(".", "")
)

type Manager struct {
	tbot *goTelegram.Bot
	repo repository.Repo
}

func (m *Manager) SendGenericMessage(text string, userId int) {
	chat := goTelegram.Chat{ID: userId}

	_, err := m.tbot.SendMessage(text, chat)

	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) SendSpendingData(userID int, spending ...model.Spending) error {
	text := strings.Builder{}

	if len(spending) == 0 {
		m.SendGenericMessage("No Spending Data Found For You", userID)
		return nil
	}
	var (
		previousTime string
		total        int
	)
	for _, record := range spending {
		tm := record.CreatedAt.Format("2006-01-02")
		if previousTime != tm {
			if previousTime != "" {
				text.WriteString(fmt.Sprintf("Total: %.2f\n\n", float64(total)/100))
				total = 0.00
			}
			text.WriteString(fmt.Sprintf("======== Date: %s ========\n", tm))
		}

		total += record.Amount
		amt := float64(record.Amount) / 100
		text.WriteString(fmt.Sprintf("Amount: %.2f\nCategory: %s\nNote: %s\n\n", amt, record.Category, record.Note))
		previousTime = tm
	}

	text.WriteString(fmt.Sprintf("Total: %.2f\n\n", float64(total)/100))

	m.SendGenericMessage(text.String(), userID)

	return nil
}

func (m *Manager) RetrieveSpendingByDateRanges(userID int, ranges ...string) ([]model.Spending, error) {
	var (
		start       string
		stop        string
		validRanges []time.Time
	)

	switch {
	case len(ranges) == 0:
		start = time.Now().Format("2006-01-02")
		stop = start
	case len(ranges) == 1:
		_, err := time.Parse("2006-01-02", ranges[0])

		if err != nil {
			m.SendGenericMessage(fmt.Sprintf("Invalid Input Entered: %v.", ranges[0]), userID)
			return nil, err
		}

		start = ranges[0]
		stop = start
	case len(ranges) == 2:
		st, err := time.Parse("2006-01-02", ranges[0])
		if err != nil {
			m.SendGenericMessage(fmt.Sprintf("Invalid Input Entered: %v.", ranges[0]), userID)
			return nil, err
		}
		et, err := time.Parse("2006-01-02", ranges[1])
		if err != nil {
			m.SendGenericMessage(fmt.Sprintf("Invalid Input Entered: %v.", ranges[0]), userID)
			return nil, err
		}

		if st.After(et) {
			m.SendGenericMessage("Start Date Cannot Be After End Date", userID)
			return nil, errors.New("start Date Is After End Date")
		}

		start = ranges[0]
		stop = ranges[1]
	default:
		m.SendGenericMessage("Received More Than Two Date Ranges, Picking The First Two", userID)

		for _, date := range ranges {
			if dt, err := time.Parse("2006-01-02", date); err == nil {
				validRanges = append(validRanges, dt)
			}
		}

		if validRanges[0].After(validRanges[1]) {
			m.SendGenericMessage("Start Date Cannot Be After End Date", userID)
			return nil, errors.New("start date is after end date")
		}

		start = validRanges[0].Format("2006-01-02")
		stop = validRanges[1].Format("2006-01-02")
	}

	spendings, err := m.repo.RetrieveSpending(userID, start, stop)

	if err != nil {
		return nil, err
	}

	return spendings, nil
}

func (m *Manager) RetrieveYesterdaySpending(userID int) ([]model.Spending, error) {
	yesterday := time.Now().AddDate(0, 0, -1)

	spendings, err := m.RetrieveSpendingByDateRanges(userID, yesterday.Format("2006-01-02"))

	if err != nil {
		return nil, err
	}

	return spendings, nil
}

func (m *Manager) RetrieveThisWeekSpending(userID int) ([]model.Spending, error) {
	now := time.Now()
	dow := int(now.Weekday())
	if dow == 0 {
		dow = 7
	}
	dow--
	monDate := now.AddDate(0, 0, -dow)

	nowString := now.Format("2006-01-02")
	monDateString := monDate.Format("2006-01-02")
	spendings, err := m.RetrieveSpendingByDateRanges(userID, monDateString, nowString)

	if err != nil {
		return nil, err
	}

	return spendings, nil
}

func (m *Manager) RetrieveThisMonthSpending(userID int) ([]model.Spending, error) {
	now := time.Now()
	//	daysInAMonth := time.Date(now.Year(), now.Month() + 1, 0, 0, 0, 0, 0, time.UTC).Day()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	today := now.Format("2006-01-02")

	spendings, err := m.RetrieveSpendingByDateRanges(userID, firstDayOfMonth, today)

	if err != nil {
		return nil, err
	}

	return spendings, nil
}

func NewManager(tbot *goTelegram.Bot, repo repository.Repo) *Manager {
	return &Manager{
		tbot: tbot,
		repo: repo,
	}
}

func (m *Manager) SaveEntry(userID, messageID int, text string) error {
	bits := strings.Fields(text)
	if !decimalRegex.MatchString(bits[0]) {
		return errors.New("invalid input detected")
	}

	date := time.Now()

	var (
		note string
	)
	amtStr := bits[0]
	if !strings.Contains(amtStr, ".") {
		amtStr = fmt.Sprintf("%s00", amtStr)
	} else {
		amtStr = dotRemover.Replace(amtStr)
	}

	if len(bits) > 1 {
		length := len(bits)
		if dateRegex.MatchString(bits[length-1]) {
			var err error
			date, err = time.Parse("2006-01-02", bits[length-1])
			if err != nil {
				return err
			}
			note = strings.Join(bits[1:length-1], " ")
		} else {
			note = strings.Join(bits[1:], " ")
		}
	}

	amt, err := strconv.Atoi(amtStr)

	if err != nil {
		return err
	}

	entry := model.Spending{
		MessageID: messageID,
		Amount:    amt,
		Note:      note,
		CreatedAt: date,
		Category:  "General",
	}

	return m.repo.SaveEntry(userID, entry)
}
