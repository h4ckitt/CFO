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
	//textRegex    = regexp.MustCompile(`([A-Za-z0-9]+).*`)
	dotRemover = strings.NewReplacer(".", "")
)

type Manager struct {
	tbot    *goTelegram.Bot
	repo    repository.Repo
	pending model.PendingReplyCache
}

func (m *Manager) SendGenericMessage(text string, userId int) {
	chat := goTelegram.Chat{ID: userId}

	_, err := m.tbot.SendMessage(text, chat)

	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) EditText(text string, messageID, chatID int) {
	message := goTelegram.Message{MessageID: messageID, Chat: goTelegram.Chat{ID: chatID}}

	_, err := m.tbot.EditMessage(message, text)

	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) DeleteMessage(messageID, chatID int) {
	message := goTelegram.Message{MessageID: messageID, Chat: goTelegram.Chat{ID: chatID}}

	err := m.tbot.DeleteMessage(message)

	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) AnswerCallBackQuery(queryID string) error {
	err := m.tbot.AnswerCallback(queryID, "", false)

	return err
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

func (m *Manager) SendCategorySpending(userID int, spending ...model.CategorySpending) error {

	text := strings.Builder{}

	if len(spending) == 0 {
		m.SendGenericMessage("No Spending Data Found For You", userID)
		return nil
	}

	var total float64

	text.WriteString("==== Spending Data By Categories ====\n")

	for _, data := range spending {
		amt := float64(data.TotalSpending) / 100
		total += amt
		text.WriteString(fmt.Sprintf("Category: %s\nAmount: %.2f\n\n", data.Category, amt))
	}

	text.WriteString(fmt.Sprintf("=====================================\nTotal: %.2f\n=====================================", total))

	m.SendGenericMessage(text.String(), userID)

	return nil
}

func (m *Manager) RetrieveCategoriesSpendingByDateRanges(userID int, ranges ...string) ([]model.CategorySpending, error) {
	var (
		start       string
		stop        string
		validRanges = make([]string, 2)
	)

	switch {
	case len(ranges) == 0:
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		stop = now.Format("2006-01-02")

	case len(ranges) == 1:
		if dateRegex.MatchString(ranges[0]) {
			start = ranges[0]
			stop = time.Now().Format("2006-01-02")
		} else {
			m.SendGenericMessage("Invalid date specified", userID)
			return nil, errors.New("invalid date specified")
		}

		st, _ := time.Parse("2006-01-02", start)
		et, _ := time.Parse("2006-01-02", stop)

		if st.After(et) {
			m.SendGenericMessage("Start Date Cannot Be After Today", userID)
			return nil, errors.New("start date after current date")
		}

	case len(ranges) == 2:
		st, err := time.Parse("2006-01-02", ranges[0])

		if err != nil {
			m.SendGenericMessage("Invalid Start Date Specified", userID)
			return nil, errors.New("invalid start date specified")
		}

		et, err := time.Parse("2006-01-02", ranges[1])

		if err != nil {
			m.SendGenericMessage("Invalid End Date Specified", userID)
			return nil, errors.New("invalid end date specified")
		}

		start = st.Format("2006-01-02")
		stop = et.Format("2006-01-02")

	default:
		numValid := 0
		for _, date := range ranges {
			if dateRegex.MatchString(date) {
				validRanges[numValid] = date
				numValid++
			}

			if numValid > 2 {
				break
			}
		}

		switch numValid {
		case 0:
			now := time.Now()
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			stop = now.Format("2006-01-02")

		case 1:
			now := time.Now()
			start = validRanges[0]
			t, _ := time.Parse("2006-01-02", start)

			if t.After(now) {
				m.SendGenericMessage("Start Date Cannot Be After Current Date", userID)
				return nil, errors.New("start date is after current date")
			}

			stop = now.Format("2006-01-02")
		default:
			start = validRanges[0]
			stop = validRanges[1]
			st, _ := time.Parse("2006-01-02", start)
			et, _ := time.Parse("2006-01-02", stop)

			if st.After(et) {
				m.SendGenericMessage("Start Date Cannot Be After Current Date", userID)
				return nil, errors.New("start date is after current date")
			}
		}
	}

	spending, err := m.repo.RetrieveSpendingByCategory(userID, start, stop)

	if err != nil {
		return nil, err
	}

	return spending, nil
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

func NewManager(tbot *goTelegram.Bot, repo repository.Repo, pending model.PendingReplyCache) *Manager {
	return &Manager{
		tbot:    tbot,
		repo:    repo,
		pending: pending,
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
	if strings.Contains(amtStr, ".") {
		amtStr = dotRemover.Replace(amtStr)
	} else {
		amtStr = fmt.Sprintf("%s00", amtStr)
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

	categories, err := m.repo.RetrieveUserCategories(userID)

	if err != nil {
		return err
	}

	if len(categories) > 0 {
		categoriesToSend := make([]string, len(categories)*2)

		for i, j := 0, 0; i < len(categories); i, j = i+1, j+2 {
			categoriesToSend[j], categoriesToSend[j+1] = categories[i], fmt.Sprintf("category-%s", categories[i])
		}

		m.tbot.CreateKeyboard(userID, 3)
		err = m.tbot.AddButtons(userID, categoriesToSend...)

		if err != nil {
			return err
		}

		m.SendGenericMessage("Select A Category Or Type A New One", userID)
	} else {
		m.SendGenericMessage("Please Enter A Category Name, It'll Be Saved Against Your Next Spend", userID)
	}

	entry := model.Spending{
		MessageID: messageID,
		Amount:    amt,
		Note:      note,
		CreatedAt: date,
		Category:  "General",
	}

	m.pending.CacheSpending(userID, &entry)

	//return m.repo.SaveEntry(userID, entry)
	return nil
}

func (m *Manager) CompleteSaveEntry(userID int, category string, shouldSave bool) error {
	cachedSpending := m.pending.RetrieveSpending(userID)

	cachedSpending.Category = category

	m.pending.DeleteCachedSpending(userID)

	if shouldSave {
		err := m.repo.SaveCategory(userID, category)

		if err != nil {
			return err
		}
	}

	m.tbot.DeleteKeyboard(userID)

	return m.repo.SaveEntry(userID, *cachedSpending)
}
