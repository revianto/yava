package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	// appEnv := os.Getenv("APP_ENV")
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv(key)
}
