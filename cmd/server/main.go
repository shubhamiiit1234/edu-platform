package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/database"
	"github.com/yourname/edu-backend-starter/internal/routes"
)

func main() {
	// 1. Connect to Database
	connection := "host=host.docker.internal port=5432 user=postgres password=mysecretpassword dbname=Edu_Platform sslmode=disable"
	err := database.InitializeDB(connection)
	if err != nil {
		fmt.Println("error initializing database: ", err.Error())
		return
	}

	// 2. Port setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3. Gin setup
	r := gin.Default()

	// --- CORS Middleware ---
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// ------------------------

	routes.RegisterRoutes(r)

	log.Println("Starting server on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
