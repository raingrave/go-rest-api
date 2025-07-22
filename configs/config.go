package configs

import (
	"fmt"
	"os"
	"strconv"
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
		return 60 // Default to 60 minutes
	}
	return minutes
}
