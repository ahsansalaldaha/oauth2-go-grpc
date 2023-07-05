package services

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBService - for generation
type DBService struct {
	DB *gorm.DB
}

// NewDBService - Generates Database Service
func NewDBService() *DBService {
	dbInstance, err := gorm.Open(mysql.New(mysql.Config{
		DSN: os.Getenv("DATABASE_URL"),
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return &DBService{
		DB: dbInstance,
	}
}
