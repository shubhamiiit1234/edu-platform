package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/routes"

	"github.com/go-chi/chi/v5"
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

	// 3. Chi Router setup (NO CORS)
	r := chi.NewRouter()

	// Register all API routes
	routes.RegisterRoutes(r)

	log.Println("Starting server on port", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
