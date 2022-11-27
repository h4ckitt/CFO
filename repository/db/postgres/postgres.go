package postgres

import (
	"cfo/config"
	"cfo/model"
	"database/sql"
	"fmt"
	"net"
	"time"

	_ "github.com/lib/pq"
)

type PgresRepo struct {
	conn *sql.DB
}

func NewPgresRepo() (*PgresRepo, error) {
	pgConfig := config.GetConfig().DB
	dbRepo := new(PgresRepo)
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user='%s' password='%s' dbname=%s sslmode=disable", pgConfig.IP, pgConfig.Port, pgConfig.UserName, pgConfig.Password, pgConfig.DBName)

	dbRepo.conn, err = sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if pgConfig.Wait {
		for !dbRepo.WaitForDB(pgConfig.IP, pgConfig.Port) {
			time.Sleep(1 * time.Second)
			continue
		}
	}

	return dbRepo, dbRepo.conn.Ping()

}

func (p PgresRepo) WaitForDB(addr, port string) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))

	return err == nil
}

func (p PgresRepo) SaveEntry(userId int, entry model.Spending) error {
	statement := `INSERT INTO spending (messageID, userID, amount, category, createdAt, note) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := p.conn.Exec(statement, entry.MessageID, userId, entry.Amount, entry.Category, entry.CreatedAt, entry.Note)
	return err
}

func (p PgresRepo) RetrieveSpending(userId int, start, end string) ([]model.Spending, error) {
	var result []model.Spending

	var (
		statement string
		args      []any
	)

	args = append(args, userId)

	if start == end {
		statement = `SELECT amount, category, createdAt, note FROM spending WHERE userID = $1 AND DATE(createdAt) = $2`
		args = append(args, start)
	} else {
		statement = `SELECT amount, category, createdAt, note FROM spending WHERE userID = $1 AND DATE(createdAt) >= $2 AND DATE(createdAt) <= $3 ORDER BY createdAt`
		args = append(args, start, end)
	}

	rows, err := p.conn.Query(statement, args...)

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

func (p PgresRepo) RetrieveSpendingByCategory(userId int, start, end string) ([]model.CategorySpending, error) {
	var result []model.CategorySpending

	var (
		statement string
		args      []any
	)

	args = append(args, userId)

	if start == end {
		statement = `SELECT SUM(amount), category FROM spending WHERE userID = $1 AND DATE(createdAt) = $2 GROUP BY category`
		args = append(args, start)
	} else {
		statement = `SELECT SUM(amount), category FROM spending WHERE userID = $1 AND DATE(createdAt) >= $2 AND DATE(createdAt) <= $3 GROUP BY category`
		args = append(args, start, end)
	}

	rows, err := p.conn.Query(statement, args...)

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

func (p PgresRepo) RetrieveUserCategories(userID int) ([]string, error) {
	var result []string

	res, err := p.conn.Query("SELECT category FROM categories WHERE userID = $1 OR userID = $2", userID, 0)

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

func (p PgresRepo) SaveCategory(userId int, category string) error {
	_, err := p.conn.Exec("INSERT INTO categories (userId, category) VALUES ($1, $2)", userId, category)

	return err
}
