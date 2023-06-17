package services

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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

var (
	queueName   = "login-logs"
	concurrency = 5
)

// MSGHandler - Handler for the message
type MSGHandler func(msg *amqp.Delivery)

// QueueService - A service for Queue Management
type QueueService struct {
	queueConn *amqp.Connection
	handler   MSGHandler
}

// NewQueueService - New Queue Service
func NewQueueService(queue *amqp.Connection, handler MSGHandler) *QueueService {
	return &QueueService{
		queueConn: queue,
		handler:   handler,
	}
}

// HandleQueueConsumer - handlers the consumer queue
func (qs *QueueService) HandleQueueConsumer() {
	logrus.Info("Running HandleQueueConsumer")

	// Create a new channel
	ch, err := qs.queueConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Set the QoS for the channel to limit the number of unacknowledged messages
	err = ch.Qos(concurrency, 0, false)
	if err != nil {
		log.Fatalf("Failed to set channel QoS: %v", err)
	}

	// Start concurrent workers to handle messages
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			qs.handleMessages(ch, workerID)
			fmt.Printf("Worker %d finished\n", workerID)
		}(i + 1)
	}

	// Wait for all workers to finish
	wg.Wait()
}

func (qs *QueueService) handleMessages(ch *amqp091.Channel, workerID int) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Worker %d: Failed to consume messages: %v", workerID, err)
		return
	}

	for msg := range msgs {
		log.Printf("Worker %d: Received message: %s", workerID, string(msg.Body))

		// Process the message here...
		go qs.handler(&msg)

		// Acknowledge the message after processing
		err := msg.Ack(false)
		if err != nil {
			log.Printf("Worker %d: Failed to acknowledge message: %v", workerID, err)
		}
	}
}
