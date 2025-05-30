package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"student-api/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	totalStudents = flag.Int("total", 10000000, "Total number of students to insert")
	batchSize     = flag.Int("batch", 1000, "Number of students to insert per batch")
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	flag.Parse()

	dsn := buildDSNFromEnv()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("Seeding student data...")
	gofakeit.Seed(time.Now().UnixNano())

	var students []model.Student

	for i := 0; i < *totalStudents; i++ {
		students = append(students, NewFakeStudent())

		if len(students) == *batchSize {
			if err := insertBatch(db, students, i+1); err != nil {
				log.Fatalf("Batch insert failed at student %d: %v", i+1, err)
			}
			students = students[:0]
		}
	}

	// Insert remaining students
	if len(students) > 0 {
		if err := insertBatch(db, students, *totalStudents); err != nil {
			log.Fatalf("Final batch insert failed: %v", err)
		}
	}

	log.Printf("Successfully inserted %d students.\n", *totalStudents)
}

// NewFakeStudent returns a random student with UUID
func NewFakeStudent() model.Student {
	return model.Student{
		ID:    uuid.New().String(),
		Name:  gofakeit.Name(),
		Age:   gofakeit.Number(1, 120),
		Email: gofakeit.Email(),
	}
}

// insertBatch inserts a batch and logs progress
func insertBatch(db *gorm.DB, students []model.Student, count int) error {
	if err := db.Create(&students).Error; err != nil {
		return err
	}
	log.Printf("Inserted %d students...", count)
	return nil
}

// buildDSNFromEnv constructs the DSN using environment variables
func buildDSNFromEnv() string {
	required := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT"}
	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required environment variable: %s", key)
		}
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}
