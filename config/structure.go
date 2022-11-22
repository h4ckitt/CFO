package config

type config struct {
	DB         PostgresDatabase
	TBotAPIKey string
	PORT       string
}
type MySQLDatabase struct {
	DBName   string
	Port     string
	Password string
	IP       string
	UserName string
	Wait     bool
}

type PostgresDatabase struct {
	DBName   string
	Port     string
	Password string
	IP       string
	UserName string
	Wait     bool
}
