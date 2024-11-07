package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Print("No config.env file found")
	}
}

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	DBHost           string
	DBPort           string
	AppPort          string
}

func New() *Config {
	return &Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		AppPort:          os.Getenv("APP_PORT"),
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword, c.DBHost, c.DBPort, c.PostgresDB)
}
