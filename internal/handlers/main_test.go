package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/middleware"
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
	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", CreateUser)
		v1.POST("/login", Login)

		authRoutes := v1.Group("/").Use(middleware.AuthMiddleware())
		{
			authRoutes.GET("/users", ListUsers)
			authRoutes.GET("/users/:id", GetUser)
			authRoutes.PUT("/users/:id", UpdateUser)
			authRoutes.DELETE("/users/:id", DeleteUser)
		}
	}
	return r
}

func clearUserTable() {
	internal.DB.MustExec("DELETE FROM users")
}