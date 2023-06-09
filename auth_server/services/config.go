package services

import (
	"fmt"
	"invento/oauth/auth_server/models"
	"time"

	"gorm.io/gorm"
)

// ConfigRememberTime - Time we are remembering configs for
const ConfigRememberTime = 1 * time.Hour

// ConfigService - Config model for config crud
type ConfigService struct {
	db       *gorm.DB
	redisSVC *RedisService
}

// NewConfigService - Creates new User Model
func NewConfigService(db *gorm.DB, redisSVC *RedisService) *ConfigService {
	return &ConfigService{
		db:       db,
		redisSVC: redisSVC,
	}
}

// GetFromDB -  get config from the database
func (configSVC *ConfigService) GetFromDB(key string) *models.Config {
	var config models.Config
	if err := configSVC.db.Where("key = ?", key).First(&config).Error; err != nil {
		return nil
	}
	return &config
}

// Get - get config
func (configSVC *ConfigService) Get(key string) (*models.Config, bool) {

	res := configSVC.redisSVC.Remember(key, ConfigRememberTime, func() interface{} {
		return configSVC.GetFromDB(key)
	})

	if res != nil {
		return res.(*models.Config), true
	}
	return nil, false

}

// FindOrNew - Find existing from key or create new config from value
func (configSVC *ConfigService) FindOrNew(key string, value interface{}) (*models.Config, error) {

	existingConfig, ok := configSVC.Get(key)
	if ok == false {
		newConfig := models.Config{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		}
		result := configSVC.db.Create(&newConfig)
		if result.Error != nil {
			// Handle the error
			return nil, fmt.Errorf("Failed to create config: %v", result.Error)
		}
		return &newConfig, nil
	}
	return existingConfig, nil

}

// Set - set configs
func (configSVC *ConfigService) Set(key string, config models.Config) (*models.Config, error) {
	configSVC.redisSVC.Set(key, config, ConfigRememberTime)
	return &config, nil
}
