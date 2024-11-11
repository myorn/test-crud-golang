package mq

import (
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"

	"event-monitor-service/events"
)

func Connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return conn
}

func ReceiverLoop(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %v", err)
	}

	defer func() {
		if err := ch.Close(); err != nil {
			log.Errorf("Failed to close channel: %v", err)
		}
	}()

	// these props below should be set on prod differently
	q, err := ch.QueueDeclare(
		"equipment", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Errorf("Failed to declare a queue: %v", err)
	}

	delivery, err := ch.Consume(
		q.Name, // routing key
		"",     // consumer
		true,   // auto-ack
		false,  // immediate
		false,  // no-local
		false,  // exclusive
		nil,    // args
	)
	if err != nil {
		log.Errorf("Failed to publish a message: %v", err)
	}

	var forever chan struct{}

	go func() {
		for d := range delivery {
			events.Router(d.Body)
		}
	}()

	log.Printf("Loop started")

	// this could be a graceful shutdown
	<-forever
}
