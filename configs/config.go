package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvDBConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func EnvJWTSecretKey() string {
	return os.Getenv("JWT_SECRET_KEY")
}

func EnvJWTExpirationMinutes() int {
	minutes, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		return 60
	}
	return minutes
}

// LoadEnvForTests finds and loads the .env.test file.
func LoadEnvForTests() {
	err := godotenv.Load("../../.env.test") // Adjust path as needed
	if err != nil {
		// Try loading from the root, in case tests are run from there
		if err = godotenv.Load(".env.test"); err != nil {
			log.Fatalf("Error loading .env.test file for tests")
		}
	}
}