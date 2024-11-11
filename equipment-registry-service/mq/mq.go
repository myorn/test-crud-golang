package mq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return conn
}
