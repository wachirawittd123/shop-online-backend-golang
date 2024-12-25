package main

import (
	"log"
	"net/http"

	"github.com/wachirawittd123/shop-online-backend-golang/common"

	"github.com/gin-gonic/gin"
	routes "github.com/wachirawittd123/shop-online-backend-golang/routes"
)

func main() {
	// Connect to the database
	db := common.ConnectDB("mongodb://localhost:27017", "shop-online")

	// Example: Use `db` to access collections
	log.Println("Database initialized:", db.Name())

	// Initialize Gin
	r := gin.Default()

	routes.RegisterRoutes(r)

	// Example route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start the server
	r.Run(":8080")
}