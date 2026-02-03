package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init opens an sqlite database at the default path database.db
func Init() {
	path := "database.db"

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	DB = db
	log.Println("Successfully connected to database")
}
