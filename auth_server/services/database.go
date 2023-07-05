package services

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBService - for generation
type DBService struct {
	DB *gorm.DB
}

// NewDBService - Generates Database Service
func NewDBService() *DBService {
	logrus.Info(os.Getenv("DATABASE_URL"))
	// dbInstance, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	dsn := os.Getenv("DATABASE_URL")
	dbInstance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return &DBService{
		DB: dbInstance,
	}
}
