package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

//Config holds the top level config
type Config struct {
	DB   *DB
	Port string
}

//DB holds DB connection config
type DB struct {
	Dialect  string
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	Charset  string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//GetConfig returns a fresh Config, including connecting to the DB
func GetConfig() *Config {
	return &Config{
		DB: &DB{
			Dialect:  "mysql",
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     os.Getenv("DB_DATABASE"),
			Charset:  "utf8",
		},
		Port: getEnv("PORT", "8080"),
	}
}

//GetURI builds a DB connection string
func (d *DB) GetURI() string {
	switch d.Dialect {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
			d.Username,
			d.Password,
			d.Host,
			d.Port,
			d.Name,
			d.Charset)
	case "sqlite3":
		return d.Name
	default:
		return "none"
	}

}
