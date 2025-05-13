package db

import (
	"fmt"
	"log"
	"os"
	"student-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&model.Student{}); err != nil {
		log.Fatal("Failed to auto-migrate Student model: ", err)
	}

	DB = db
	log.Println("Connected to PostgreSQL database")
}
