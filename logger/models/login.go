package models

import (
	"time"

	"gorm.io/gorm"
)

// User - Represents User Model
type User struct {
	gorm.Model
	Name     string
	Username string
	Password string
}

// LogMessage - Represents Logs Message
type LogMessage struct {
	gorm.Model
	UserID    uint
	User      User
	AttemptAt time.Time
	Success   bool
}
