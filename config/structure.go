package config

type config struct {
	DB         MySQLDatabase
	TBotAPIKey string
	PORT       string
}
type MySQLDatabase struct {
	DBName   string
	Port     string
	Password string
	IP       string
	UserName string
}
