package models

import (
	"fmt"
	"invento/oauth/auth_server/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// User - User models
type User struct {
	gorm.Model
	Name     string `validate:"required"`
	Username string `gorm:"unique" validate:"required"`
	Password string `validate:"required"`
	Lock     *[]UserLock
}

// UserModel - User model to generate clients
type UserModel struct {
	db *gorm.DB
}

// NewUserModel - Creates new User Model
func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

// GenerateUser - Generates User of certain type
func (cm *UserModel) GenerateUser(name string, username string, password string) (*User, error) {

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		panic("Failed to hash password")
	}

	user := User{Name: name, Username: username, Password: hashedPassword}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		return nil, fmt.Errorf("validation Error")
	}
	result := cm.db.Create(&user)
	if result.Error != nil {
		// Handle the error
		return nil, fmt.Errorf("Failed to create user: %v", result.Error)
	}
	user.Password = ""

	return &user, nil

}

// UserLock - User Lock Model models
type UserLock struct {
	gorm.Model
	UserRelation
	LockedAt time.Time
}

// UserRelation - represents the user relationship
type UserRelation struct {
	UserID uint  `json:"user_id" validate:"required"`
	User   *User `json:"user"`
}
