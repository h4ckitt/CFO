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
	dsn := fmt.Sprintf("%s:%s@%s(%s;%s)/%s", conf.UserName, conf.Password, "tcp", conf.IP, conf.Port, conf.DBName)
	dbRepo.conn, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return dbRepo, nil
}

func (m *MySqlRepo) SaveEntry(userId int, entry *model.Spending) error {
	return nil
}

func (m *MySqlRepo) RetrieveSpending(userId int, start, end string) ([]model.Spending, error) {
	return []model.Spending{}, nil
}

func (m *MySqlRepo) RetrieveSpendingByCategory(userId int, category string) ([]model.CategorySpending, error) {
	return []model.CategorySpending{}, nil
}
