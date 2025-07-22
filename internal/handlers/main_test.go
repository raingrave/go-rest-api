package handlers

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/raingrave/apirest/internal"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup database connection
	// Note: This will use the same DB as dev, but we will truncate tables.
	// For a real-world scenario, consider a separate test database.
	internal.ConnectDB()

	// Setup router
	router = setupRouter()

	// Run tests
	code := m.Run()

	// Exit
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
