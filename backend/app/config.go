package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB   *DBConfig
	Port int
}

type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
	URI      string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_DATABASE"),
			Charset:  "utf8",
		},
		Port: 4000,
	}
	config.DB.URI = fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	return &config
}
