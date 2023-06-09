package models

import (
	"gorm.io/gorm"
)

// Config - Config models
type Config struct {
	gorm.Model
	Key   string
	Value string
}
