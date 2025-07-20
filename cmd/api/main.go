package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raingrave/apirest/internal"
	"github.com/raingrave/apirest/internal/handlers"
)

func main() {
	internal.ConnectDB()

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/users", handlers.ListUsers)
	r.POST("/users", handlers.CreateUser)

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/:id", handlers.GetUser)
		userRoutes.PUT("/:id", handlers.UpdateUser)
		userRoutes.DELETE("/:id", handlers.DeleteUser)
	}

	r.Run(":3000")
}
