package services

import (
	"fmt"
	"invento/oauth/auth_server/models"
	"invento/oauth/auth_server/utils"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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

// GetDatabaseService - returns the database service assigned in config service
func (configSVC *ConfigService) GetDatabaseService() *gorm.DB {
	return configSVC.db
}

// GetFromDB -  get config from the database
func (configSVC *ConfigService) GetFromDB(key string) *models.Config {
	var config models.Config
	if err := configSVC.db.Where("`key` = ?", key).First(&config).Error; err != nil {
		logrus.Infof("Error found in get from db: ", err)
		return nil
	}
	logrus.Infof("config got from DB: %v", config)
	return &config
}

// Get - get config
func (configSVC *ConfigService) Get(key string) (*models.Config, bool) {
	res := configSVC.redisSVC.Remember(key, utils.ConfigRememberTime, func() interface{} {
		logrus.Info("get from db called")
		return configSVC.GetFromDB(key)
	})
	// logrus.Info(res)
	// logrus.Infof("res found %v", res)

	if res != nil {
		// logrus.Info("A1")
		// logrus.Infof("%v", res)
		return res.(*models.Config), true
	}
	// logrus.Info("A2")
	return nil, false
}

// GetBool - get config and parse to bool
func (configSVC *ConfigService) GetBool(key string) (bool, bool) {
	if obj, ok := configSVC.Get(key); ok == true {
		b, err := strconv.ParseBool(obj.Value)
		if err != nil {
			fmt.Println("Error parsing boolean:", err)
		}
		return b, true
	}
	return false, false
}

// FindOrCreate - Find existing from key or create new config from value
func (configSVC *ConfigService) FindOrCreate(key string, value interface{}) (*models.Config, error) {
	existingConfig, ok := configSVC.Get(key)
	// logrus.Infof("existingConfig %v", existingConfig)
	// logrus.Infof("ok %v", ok)

	if ok == false {
		newConfig := models.Config{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		}
		result := configSVC.db.Create(&newConfig)
		if result.Error != nil {
			// logrus.Info("A")
			// Handle the error
			return nil, fmt.Errorf("Failed to create config: %v", result.Error)
		}
		// logrus.Info("B")
		return &newConfig, nil
	}
	// logrus.Info("C")
	if existingConfig == nil {
		newConfig := models.Config{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		}
		result := configSVC.db.Create(&newConfig)
		if result.Error != nil {
			// logrus.Info("CA")
			// Handle the error
			return nil, fmt.Errorf("Failed to create config: %v", result.Error)
		}
		// logrus.Info("CB")
		return &newConfig, nil
	}
	return existingConfig, nil

}

// Set - set configs
func (configSVC *ConfigService) Set(key string, config models.Config) (*models.Config, error) {
	configSVC.redisSVC.Set(key, config, utils.ConfigRememberTime)
	return &config, nil
}

// GetPasswordComplexityConfig -  return Password Complexity Object
func (configSVC *ConfigService) GetPasswordComplexityConfig() utils.PasswordComplexity {
	pc := utils.PasswordComplexity{}
	if upper, ok := configSVC.GetBool("password-complexity-should-have-upper"); ok == true {
		pc.ShouldHaveUppercase = upper
	}
	if lower, ok := configSVC.GetBool("password-complexity-should-have-lower"); ok == true {
		pc.ShouldHaveLowercase = lower
	}
	if numeric, ok := configSVC.GetBool("password-complexity-should-have-numeric"); ok == true {
		pc.ShouldHaveNumber = numeric
	}
	if special, ok := configSVC.GetBool("password-complexity-should-have-special"); ok == true {
		pc.ShouldHaveSpecial = special
	}
	return pc
}

// StoreUtmostComplexity - stored utmost complexity
func (configSVC *ConfigService) StoreUtmostComplexity() {
	configSVC.FindOrCreate("password-min-length", 8)
	configSVC.FindOrCreate("password-complexity-should-have-upper", true)
	configSVC.FindOrCreate("password-complexity-should-have-lower", true)
	configSVC.FindOrCreate("password-complexity-should-have-numeric", true)
	configSVC.FindOrCreate("password-complexity-should-have-special", true)
}

// StoreDefaultConfigs - Stores Default Configs
func (configSVC *ConfigService) StoreDefaultConfigs() {
	configSVC.StoreUtmostComplexity()
	configSVC.FindOrCreate("enable-login-activity-logging", true)
}
