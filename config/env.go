package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	PG_HOST           string
	PG_PORT           string
	PG_USER           string
	PG_PASSWORD       string
	PG_DB             string
	STARK_PAY_API_KEY string
)

func LoadEnv() error {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	// Postgres details
	PG_HOST = os.Getenv("DB_HOST")
	PG_PORT = os.Getenv("DB_PORT")
	PG_USER = os.Getenv("DB_USERNAME")
	PG_PASSWORD = os.Getenv("DB_PASSWORD")
	PG_DB = os.Getenv("DB_NAME")
	// Stark Pay API key
	STARK_PAY_API_KEY = os.Getenv("STARK_PAY_API_KEY")

	return nil
}
