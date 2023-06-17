package services

import (
	"invento/oauth/auth_server/models"
	"invento/oauth/auth_server/utils"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// CredentialService - for generation
type CredentialService struct {
	dbSVC     *DBService
	redisSVC  *RedisService
	queueSVC  *QueueService
	lockCount map[string]int
}

// ValidateClientAndSecret - validates for the client and secret specified are correct
func (credSVC *CredentialService) ValidateClientAndSecret(clientID string, secret string) bool {
	logrus.Info("ValidateClientAndSecret: ", clientID, secret)
	res := credSVC.getRememberedClient(clientID)
	if res != nil && res.(models.Client).Secret == secret {
		return true
	}
	return false
}

func (credSVC *CredentialService) getRememberedClient(clientID string) interface{} {
	return credSVC.redisSVC.Remember(clientID, utils.CredsClientCache, func() interface{} {
		var client models.Client
		if err := credSVC.dbSVC.DB.Where("client_id = ?", clientID).Preload("Redirects").First(&client).Error; err != nil {
			return nil
		}
		return client
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
	logrus.Info("username: ", username)
	logrus.Info("password: ", password)
	var user models.User
	if err := credSVC.dbSVC.DB.Select("users.*").Where("username = ?", username).Joins("LEFT JOIN user_locks ON users.id = user_locks.user_id").Where("user_locks.user_id IS NULL").First(&user).Error; err != nil {
		logrus.Info("user not found")
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		credSVC.lockCount[username]++
		logrus.Infof("Lock Count %s: %d", username, credSVC.lockCount[username])
		if credSVC.lockCount[username] >= utils.UserIncorrectLoginCount {
			userLock := models.UserLock{}
			userLock.LockedAt = time.Now()
			userLock.UserID = user.ID
			credSVC.dbSVC.DB.Save(&userLock)
			logrus.Infof("Locked the profile: %s", username)
			credSVC.lockCount[username] = 0
		}
		logrus.Infof("Sending False: User found: %v", user)
		credSVC.queueSVC.ProduceMsg(LogMessage{
			User:      user,
			AttemptAt: time.Now(),
			Success:   false,
		})
	} else {
		logrus.Info("User/Password found")
		credSVC.lockCount[username] = 0
		logrus.Infof("Sending True: User found: %v", user)
		credSVC.queueSVC.ProduceMsg(LogMessage{
			User:      user,
			AttemptAt: time.Now(),
			Success:   true,
		})
	}

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
func NewCredentialService(db *DBService, rs *RedisService, queue *QueueService) *CredentialService {
	return &CredentialService{
		dbSVC:     db,
		redisSVC:  rs,
		queueSVC:  queue,
		lockCount: make(map[string]int),
	}
}
