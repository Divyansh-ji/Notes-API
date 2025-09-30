package intializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	// Load .env only if running locally (not inside Docker)
	_ = godotenv.Load()

	dsn := os.Getenv("DB_url")
	if dsn == "" {
		log.Fatal("DB_url is not set in environment variables")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to database")
}
