package mysql

import (
	"cfo/config"
	"cfo/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlRepo struct {
	conn *sql.DB
}

func NewMySQLHandler() (*MySqlRepo, error) {
	conf := config.GetConfig().DB
	dbRepo := new(MySqlRepo)
	var err error
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", conf.UserName, conf.Password, "tcp", conf.IP, conf.Port, conf.DBName)
	dbRepo.conn, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return dbRepo, dbRepo.conn.Ping()
}

func (m *MySqlRepo) SaveEntry(userId int, entry model.Spending) error {
	statement := `INSERT INTO spending (messageID, userID, amount, category, createdAt, note) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.conn.Exec(statement, entry.MessageID, userId, entry.Amount, entry.Category, entry.CreatedAt, entry.Note)
	return err
}

func (m *MySqlRepo) RetrieveSpending(userId int, start, end string) ([]model.Spending, error) {
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
		statement = `SELECT amount, category, createdAt, note FROM spending WHERE userID = ? AND createdAt >= ? AND createdAt <= ?`
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

func (m *MySqlRepo) RetrieveSpendingByCategory(userId int, category string) ([]model.CategorySpending, error) {
	var result []model.CategorySpending

	statement := `SELECT SUM(amount), category FROM spending WHERE userID = ? GROUP BY category`

	rows, err := m.conn.Query(statement, userId)

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
