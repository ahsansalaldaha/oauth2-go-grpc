package services

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBService - for generation
type DBService struct {
	DB *gorm.DB
}

// NewDBService - Generates Database Service
func NewDBService() *DBService {
	dbInstance, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DATABASE_URL"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return &DBService{
		DB: dbInstance,
	}
}
