package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GrantType - type of grants available
type GrantType byte

const (
	// Code - grand type having codes
	Code GrantType = iota
	// Implicit - grand type with username password
	Implicit
)

// Client - client model for client database table
type Client struct {
	gorm.Model
	ClientID  string     `gorm:"unique" validate:"required"`
	Secret    string     `validate:"required"`
	GrantType GrantType  `validate:"required"`
	Redirects []Redirect `validate:"dive"`
}

// ClientRelation - represents the client relationship
type ClientRelation struct {
	ClientID uint    `json:"client_id" validate:"required"`
	Client   *Client `json:"client"`
}

// Redirect - Store all possible redirect urls
type Redirect struct {
	gorm.Model
	ClientRelation
	RedirectURL string `validate:"required,url"`
}

// ClientModel - Client model to generate clients
type ClientModel struct {
	db *gorm.DB
}

// NewClientModel - Creates new Client Model
func NewClientModel(db *gorm.DB) *ClientModel {
	return &ClientModel{
		db: db,
	}
}

// GenerateClient - Generates Client of certain type
func (cm *ClientModel) GenerateClient(grantType GrantType, redirect []string) (*Client, error) {

	client := Client{ClientID: uuid.New().String(), Secret: uuid.New().String(), GrantType: grantType}

	validate := validator.New()
	err := validate.Struct(client)
	if err != nil {
		return nil, fmt.Errorf("validation Error")
	}

	result := cm.db.Create(&client)
	if result.Error != nil {
		// Handle the error
		return nil, fmt.Errorf("Failed to create client: %v", result.Error)
	}

	for _, v := range redirect {
		cm.db.Create(&Redirect{
			ClientRelation: ClientRelation{
				ClientID: client.ID,
			},
			RedirectURL: v,
		})
	}

	cm.db.Model(&client).Association("Redirects").Find(&client.Redirects)

	return &client, nil

}
