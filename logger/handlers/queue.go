package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"invento/oauth/logger/models"
	"invento/oauth/logger/services"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	queueName   = "login-logs"
	concurrency = 5
)

type UserForm struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Name      string
	Username  string
	Password  string
}

func (uf *UserForm) toUser() models.User {
	user := models.User{
		Model: gorm.Model{
			ID:        uf.ID,
			CreatedAt: uf.CreatedAt,
			UpdatedAt: uf.UpdatedAt,
			DeletedAt: uf.DeletedAt,
		},
		Name:     uf.Name,
		Username: uf.Username,
		Password: uf.Password,
	}
	// logrus.Infof("User: %v", user)
	return user
}

type LogMessageForm struct {
	User      UserForm
	AttemptAt time.Time
	Success   bool
}

func (lmf *LogMessageForm) toLogMessage() models.LogMessage {
	logMessage := models.LogMessage{
		AttemptAt: lmf.AttemptAt,
		Success:   lmf.Success,
		UserID:    lmf.User.ID,
		User:      lmf.User.toUser(),
	}
	// utils.PrintStructFields(logMessage)
	return logMessage
}

// MSGHandler - handlers the MSGs
func MSGHandler(dbService *services.DBService, mailService *services.EmailService) services.MSGHandler {

	return func(msg *amqp.Delivery) {
		var logMessageForm LogMessageForm

		if err := json.Unmarshal(msg.Body, &logMessageForm); err != nil {
			logrus.Info("error parsing message", err)
		}

		logMessage := logMessageForm.toLogMessage()
		dbService.DB.Save(&logMessage)

		email := gomail.NewMessage()
		email.SetHeader("From", "ahsan@geekinn.co")
		email.SetHeader("To", "abc@invento.sa")
		email.SetHeader("Subject", "Login")
		email.SetBody("text/plain", fmt.Sprintf("This is to notify that %s tried to access %s", logMessage.User.Name, logMessage.AttemptAt))

		mailService.Send(email)
	}

}

// HandleQueueConsumer - handlers the consumer queue
func HandleQueueConsumer() {

	conn := services.ConnectQueue()
	defer services.DisconnectQueue(conn)

	dbService := services.NewDBService()
	dbService.DB.AutoMigrate(&models.LogMessage{})

	mailService := services.NewEmailService()
	queueSVC := services.NewQueueService(conn, MSGHandler(dbService, mailService))
	queueSVC.HandleQueueConsumer()

}
