package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nickysemenza/food/backend/app/model"
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
	Host     string
	Port     string
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
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     dbName,
			Charset:  "utf8",
		},
		Port: getEnv("PORT", "8080"),
	}
	return &config
}

//GetDBURI builds a DB connection string
func (config *Config) GetDBURI() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)
}
