package database

import (
	"clockit/backend/models"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(path string) error {
	if path == "" {
		return fmt.Errorf("empty database path")
	}

	// Ensure the directory for the database file exists
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Printf("Failed to create directory %q: %v", dir, err)
			return err
		}
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to sqlite database at %q: %v", path, err)
		return err
	}

	DB = db
	log.Printf("Successfully connected to database at %q", path)

	// Auto-migrate models
	if err := DB.AutoMigrate(
		&models.Employee{},
	); err != nil {
		return fmt.Errorf("auto-migrate failed: %v", err)
	}
	return nil
}
