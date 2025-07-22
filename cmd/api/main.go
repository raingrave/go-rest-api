package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/raingrave/apirest/configs"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/handlers"
	"github.com/raingrave/apirest/internal/middleware"
	"github.com/raingrave/apirest/internal/repositories"
	"github.com/raingrave/apirest/internal/services"
)

func main() {
	internal.ConnectDB()

	m, err := migrate.New(
		"file://internal/database/migrations",
		configs.EnvDatabaseURL(),
	)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not run migrations: %v", err)
	}
	log.Println("Database migrations ran successfully")

	// Dependency Injection
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService)

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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

	r.Run(":3000")
}