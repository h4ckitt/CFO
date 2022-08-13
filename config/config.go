package config

import (
	"github.com/joho/godotenv"
	"os"
)

var conf config

func ReadConfig(filename string) error {
	err := godotenv.Load(filename)

	if err != nil {
		return err
	}

	conf = config{
		DB: MySQLDatabase{
			DBName:   os.Getenv("DB_NAME"),
			UserName: os.Getenv("DB_USER_NAME"),
			Password: os.Getenv("DB_PASS"),
			Port:     os.Getenv("DB_PORT"),
			IP:       os.Getenv("DB_IP"),
		},
		TBotAPIKey: os.Getenv("TBOT_API_KEY"),
		PORT:       os.Getenv("PORT"),
	}

	return nil
}

func GetConfig() config {
	return conf
}
