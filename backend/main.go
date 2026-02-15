package main

import (
	"clockit/backend/database"
	"clockit/backend/routes"
	"log"
	"os"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Determine database path from environment variable or use default
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: could not load .env file: %v (continuing with defaults)", err)
	}

	dbPath := os.Getenv("DB_PATH")
	defaultPath := "database/database.db"

	if dbPath == "" {
		dbPath = defaultPath
	}

	// Connect to the database using GORM
	if err := database.Init(dbPath); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r) // Configure HTTP server with proper timeouts
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10 * time.Second, // max time to read request body
		ReadHeaderTimeout: 5 * time.Second,  // max time to read headers
		WriteTimeout:      15 * time.Second, // max time to write response
		IdleTimeout:       60 * time.Second, // max time for keep-alive connections
	}

	log.Println("Starting ClockIt server on :8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
