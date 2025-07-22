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
)

func main() {
	internal.ConnectDB()

	// Run database migrations
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

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", handlers.CreateUser)
		v1.POST("/login", handlers.Login)

		authRoutes := v1.Group("/").Use(middleware.AuthMiddleware())
		{
			authRoutes.GET("/users", handlers.ListUsers)
			authRoutes.GET("/users/:id", handlers.GetUser)
			authRoutes.PUT("/users/:id", handlers.UpdateUser)
			authRoutes.DELETE("/users/:id", handlers.DeleteUser)
		}
	}

	r.Run(":3000")
}
