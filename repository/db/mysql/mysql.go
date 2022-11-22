package mysql

import (
	"cfo/config"
	"cfo/model"
	"database/sql"
	"fmt"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MeSqlRepo struct {
	conn *sql.DB
}

func NewMySQLHandler() (*MeSqlRepo, error) {
	conf := config.GetConfig().DB
	dbRepo := new(MeSqlRepo)
	var err error
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", conf.UserName, conf.Password, "tcp", conf.IP, conf.Port, conf.DBName)
	dbRepo.conn, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if conf.Wait {
		for !dbRepo.WaitForDB(conf.IP, conf.Port) {
			time.Sleep(1)
			continue
		}

	}
	return dbRepo, dbRepo.conn.Ping()
}

func (m *MeSqlRepo) WaitForDB(addr, port string) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))

	if err != nil {
		return false
	}
	return true
}

func (m *MeSqlRepo) SaveEntry(userId int, entry model.Spending) error {
	statement := `INSERT INTO spending (messageID, userID, amount, category, createdAt, note) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.conn.Exec(statement, entry.MessageID, userId, entry.Amount, entry.Category, entry.CreatedAt, entry.Note)
	return err
}

func (m *MeSqlRepo) RetrieveSpending(userId int, start, end string) ([]model.Spending, error) {
	var result []model.Spending

	var (
		statement string
		args      []any
	)

	args = append(args, userId)

	if start == end {
		statement = `SELECT amount, category, createdAt, note FROM spending WHERE userID = ? AND DATE(createdAt) = ?`
		args = append(args, start)
	} else {
		statement = `SELECT amount, category, createdAt, note FROM spending WHERE userID = ? AND DATE(createdAt) >= ? AND DATE(createdAt) <= ?`
		args = append(args, start, end)
	}

	rows, err := m.conn.Query(statement, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res model.Spending

		if err := rows.Scan(&res.Amount, &res.Category, &res.CreatedAt, &res.Note); err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	return result, nil
}

func (m *MeSqlRepo) RetrieveSpendingByCategory(userId int, start, end string) ([]model.CategorySpending, error) {
	var result []model.CategorySpending

	var (
		statement string
		args      []any
	)

	args = append(args, userId)

	if start == end {
		statement = `SELECT SUM(amount), category FROM spending WHERE userID = ? AND DATE(createdAt) = ? GROUP BY category`
		args = append(args, start)
	} else {
		statement = `SELECT SUM(amount), category FROM spending WHERE userID = ? AND DATE(createdAt) >= ? AND DATE(createdAt) <= ? GROUP BY category`
		args = append(args, start, end)
	}

	rows, err := m.conn.Query(statement, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res model.CategorySpending

		if err := rows.Scan(&res.TotalSpending, &res.Category); err != nil {
			return nil, err
		}

		result = append(result, res)
	}
	return result, nil
}

func (m *MeSqlRepo) RetrieveUserCategories(userID int) ([]string, error) {
	var result []string

	res, err := m.conn.Query("SELECT category FROM categories WHERE userID = ? OR userID = ?", userID, 0)

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var category string

		if err = res.Scan(&category); err != nil {
			return nil, err
		}

		result = append(result, category)
	}

	return result, nil
}

func (m *MeSqlRepo) SaveCategory(userId int, category string) error {
	_, err := m.conn.Exec("INSERT INTO categories (userId, category) VALUES (?, ?)", userId, category)

	return err
}
