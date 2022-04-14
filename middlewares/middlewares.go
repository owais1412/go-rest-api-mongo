package middlewares

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DotEnvVariable(key string) string {

	// for testing we need to pass the env path as like this
	// err := godotenv.Load("../.env")
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func IsValidToken(token string) bool {
	return token == "owais"
}
