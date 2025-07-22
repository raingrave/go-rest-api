package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/raingrave/apirest/internal"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env.test file for tests: %v", err)
	}

	internal.ConnectDB()

	router = setupRouter()

	code := m.Run()

	os.Exit(code)
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