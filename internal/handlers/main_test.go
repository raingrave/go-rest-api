package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/middleware"
	"github.com/raingrave/apirest/internal/repositories"
	"github.com/raingrave/apirest/internal/services"
)

var router *gin.Engine
var userRepo repositories.UserRepository

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env.test file for tests: %v", err)
	}

	internal.ConnectDB()
	userRepo = repositories.NewUserRepository()

	router = setupRouter()

	code := m.Run()
	os.Exit(code)
}

func setupRouter() *gin.Engine {
	r := gin.New()

	// DI for tests
	userService := services.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)
	authHandler := NewAuthHandler(userService)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
		v1.POST("/login", authHandler.Login)

		authRoutes := v1.Group("/").Use(middleware.AuthMiddleware())
		{
			authRoutes.GET("/users", userHandler.ListUsers)
			authRoutes.GET("/users/:id", userHandler.GetUser)
			authRoutes.PUT("/users/:id", userHandler.UpdateUser)
			authRoutes.DELETE("/users/:id", userHandler.DeleteUser)
		}
	}
	return r
}

func clearUserTable() {
	internal.DB.MustExec("DELETE FROM users")
}
