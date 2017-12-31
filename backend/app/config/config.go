package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB   *DBConfig
	Port string
}

type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
}

type Env struct {
	DB   *gorm.DB
	Port string
	Host string
	//Router **mux.Router
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	dbName := ""
	if _, ok := os.LookupEnv("TESTMODE"); ok {
		dbName = os.Getenv("DB_DATABASE_TEST")
		log.Printf("TESTMODE! using db: %s", dbName)
	} else {
		dbName = os.Getenv("DB_DATABASE")
		log.Printf("using db: %s", dbName)
	}

	config := Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     dbName,
			Charset:  "utf8",
		},
		Port: getEnv("PORT", "4000"),
	}
	return &config
}

func (config *Config) GetDBURI() string {
	return fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)
}
