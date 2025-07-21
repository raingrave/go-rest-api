package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/handlers"
	"github.com/raingrave/apirest/internal/middleware"
)

func main() {
	internal.ConnectDB()

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint (unversioned)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Versioned API group
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.POST("/users", handlers.CreateUser)
		v1.POST("/login", handlers.Login)

		// Protected routes
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
