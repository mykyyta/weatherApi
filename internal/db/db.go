package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"weatherApi/internal/model"
)

var DB *gorm.DB

func InitDatabase(dbType, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	if dbType == "postgres" {
		err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error
		if err != nil {
			return nil, fmt.Errorf("failed to enable pgcrypto: %v", err)
		}
	}

	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}

	return db, nil
}

func ConnectDefaultDB() {
	dbType := os.Getenv("DB_TYPE") // "postgres"
	dsn := os.Getenv("DB_URL")

	var err error
	DB, err = InitDatabase(dbType, dsn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB and ran AutoMigrate")
}
