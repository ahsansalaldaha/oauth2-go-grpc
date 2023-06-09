package services

import (
	"invento/oauth/auth_server/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// CredentialService - for generation
type CredentialService struct {
	dbSVC    *DBService
	redisSVC *RedisService
}

// ValidateClientAndSecret - validates for the client and secret specified are correct
func (credSVC *CredentialService) ValidateClientAndSecret(clientID string, secret string) bool {
	res := credSVC.getRememberedClient(clientID)
	if res != nil && res == secret {
		return true
	}
	return false
}

func (credSVC *CredentialService) getRememberedClient(clientID string) interface{} {
	return credSVC.redisSVC.Remember(clientID, 1*time.Hour, func() interface{} {
		var client models.Client
		if err := credSVC.dbSVC.DB.Where("client_id = ?", clientID).Preload("Redirects").First(&client).Error; err != nil {
			return nil
		}
		return client.Secret
	})
}

// ValidateClientID - validates if clientID exists
func (credSVC *CredentialService) ValidateClientID(clientID string) bool {
	res := credSVC.getRememberedClient(clientID)
	if res != nil {
		return true
	}
	return false
}

// ValidateUserCredentials - validates for the username and password specified are correct
func (credSVC *CredentialService) ValidateUserCredentials(username string, password string) bool {
	var user models.User
	if err := credSVC.dbSVC.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ValidateRedirectURI - validates Redirect Url against client ID
func (credSVC *CredentialService) ValidateRedirectURI(clientID string, redirectURI string) bool {
	res := credSVC.getRememberedClient(clientID)
	if res != nil {
		for _, redirect := range res.(models.Client).Redirects {
			if redirect.RedirectURL == redirectURI {
				return true
			}
		}
	}
	return false
}

// NewCredentialService -
func NewCredentialService(db *DBService, rs *RedisService) *CredentialService {
	return &CredentialService{
		dbSVC:    db,
		redisSVC: rs,
	}
}
