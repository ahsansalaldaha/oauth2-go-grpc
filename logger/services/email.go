package services

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// EmailService - A Service for Email Management
type EmailService struct {
	mailer *gomail.Dialer
}

// Send - Send Email message
func (es *EmailService) Send(email *gomail.Message) {
	// Send the email
	if err := es.mailer.DialAndSend(email); err != nil {
		logrus.Println("Failed to send email:", err)
		return
	}
}

// NewEmailService - Create new Email Sending Service
func NewEmailService() *EmailService {
	// Create a new SMTP sender
	// d := gomail.NewDialer("smtp.example.com", 587, "your_email@example.com", "your_password")
	port, err := strconv.Atoi(os.Getenv("MAILER_PORT"))
	if err != nil {
		logrus.Println("Error:", err)
		panic("Error converting mailer port to int")
	}

	d := gomail.NewDialer(os.Getenv("MAILER_HOST"), port, os.Getenv("MAILER_USERNAME"), os.Getenv("MAILER_PASSWORD"))
	// "mailhog", 1025, "", ""

	return &EmailService{
		mailer: d,
	}
}
