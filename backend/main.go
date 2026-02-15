package main

import (
	"clockit/backend/database"
	"clockit/backend/routes"
	"log"
	"os"

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
	routes.RegisterRoutes(r)

	log.Println("Starting ClockIt server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
