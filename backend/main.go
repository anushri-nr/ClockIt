package main

import (
	"clockit/backend/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "OK")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	// Determine database path from environment variable or use default
	defaultPath := "database/database.db"
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultPath
	}

	// Connect to the database using GORM
	if err := database.Init(dbPath); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// Set timeouts to avoid Slowloris attacks.
		// Tune these values as needed.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
