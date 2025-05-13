package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"weatherApi/internal/model"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error
	if err != nil {
		log.Fatalf("failed to enable pgcrypto: %v", err)
	}

	err = DB.AutoMigrate(&model.Subscription{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	fmt.Println("Connected to DB and ran AutoMigrate")
}
