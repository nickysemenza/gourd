package config

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/nickysemenza/food/backend/app/model"
	"log"
	"os"
)

//Config holds the top level config
type Config struct {
	DB   *DBConfig
	Port string
}

//DBConfig holds DB connection config
type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
}

//Env holds misc env stuff like the DB connection object.
type Env struct {
	DB          *gorm.DB
	Port        string
	Host        string
	CurrentUser *model.User
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//GetConfig returns a fresh Config, including connecting to the DB
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

//GetDBURI builds a DB connection string
func (config *Config) GetDBURI() string {
	return fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)
}
