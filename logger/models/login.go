package models

import (
	"time"

	"gorm.io/gorm"
)

// User - Represents User Model
type User struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	Username string `gorm:"size:255"`
	Password string `gorm:"size:255"`
}

// LogMessage - Represents Logs Message
type LogMessage struct {
	gorm.Model
	UserID    uint
	User      User
	AttemptAt time.Time
	Success   bool
}
