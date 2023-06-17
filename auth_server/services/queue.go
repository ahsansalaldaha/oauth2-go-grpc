package services

import (
	"encoding/json"
	"invento/oauth/auth_server/models"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	// QueueName signifies logging channel
	QueueName = "login-logs"
)

// ConnectQueue ...
func ConnectQueue() *amqp.Connection {
	queueConn, err := amqp.Dial(os.Getenv("QUEUE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return queueConn
}

// DisconnectQueue - Disconnects the queue
func DisconnectQueue(queueConn *amqp.Connection) {
	defer queueConn.Close()
}

// QueueService - service to handle all queue features
type QueueService struct {
	conn      *amqp.Connection
	queueName string
}

// NewQueueService - represents queue service
func NewQueueService() *QueueService {
	return &QueueService{
		conn:      ConnectQueue(),
		queueName: QueueName,
	}
}

// LogMessage - Represents Logs Message
type LogMessage struct {
	User      models.User
	AttemptAt time.Time
	Success   bool
}

// ProduceMsg produces the message on queue
func (qs *QueueService) ProduceMsg(msg LogMessage) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal JSON data: %v", err)
	}

	// Create a channel
	ch, err := qs.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		qs.queueName, // queue name
		true,         // durable
		false,        // auto-delete
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.Publish(
		"",           // exchange
		qs.queueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Message sent: %s", jsonData)
}
