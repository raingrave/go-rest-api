package handlers

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/raingrave/apirest/internal"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Find and load .env file
	rootDir, err := findProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env.test file for tests")
	}

	internal.ConnectDB()
	router = setupRouter()
	code := m.Run()
	os.Exit(code)
}

// findProjectRoot finds the project root by looking for the go.mod file.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.POST("/api/v1/users", CreateUser)
	r.POST("/api/v1/login", Login)
	// ... add other routes as you write tests for them
	return r
}

func clearUserTable() {
	internal.DB.MustExec("DELETE FROM users")
}
