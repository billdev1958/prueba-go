package main

import (
	"log"
	"net/http"
	v1 "prueba-go/internal/infrastructure/http/v1"

	"github.com/gin-gonic/gin"
)

// ... (existing global annotations)

func main() {
	r := gin.Default()

	// In a real scenario, we'd initialize repos, usecases, and handlers here.
	// For now, we provide empty config to RegisterRoutes for documentation purposes.
	v1.RegisterRoutes(r, v1.RouterConfig{})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
