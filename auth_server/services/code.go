package services

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CodeService - for generation
type CodeService struct {
	rs       *RedisService
	validity time.Duration
}

// Get - get random code
func (cs *CodeService) Get() string {
	code := uuid.New().String()
	cs.rs.Set(code, code, cs.validity)
	return code
}

// Verify - verify generated code
func (cs *CodeService) Verify(code string) bool {
	logrus.Info("Code: ", code, ":Verify")
	_, err := cs.rs.Get(code).Result()
	if err == redis.Nil {
		logrus.Info("Found Nil")
		return false
	}
	logrus.Info("Found True")
	return true
}

// NewCodeService - Generates New Code Service
func NewCodeService(rs *RedisService) *CodeService {
	return &CodeService{
		rs:       rs,
		validity: 24 * time.Hour,
	}
}
