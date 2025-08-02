package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadDbConfig() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("DB_DSN")
}
