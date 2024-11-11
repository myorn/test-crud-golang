package mq

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"

	"equipment-registry-service/dto"
)

const (
	// EquipmentCreated event
	EquipmentCreated = "EquipmentCreated"

	// EquipmentRead event
	EquipmentRead = "EquipmentRead"

	// EquipmentUpdated event
	EquipmentUpdated = "EquipmentUpdated"

	// EquipmentDeleted event
	EquipmentDeleted = "EquipmentDeleted"

	// EquipmentRestored event
	EquipmentRestored = "EquipmentRestored"
)

type equipmentEvent struct {
	EventType string `json:"eventType"`

	dto.Equipment
}

func SendEvent(ctx context.Context, conn *amqp.Connection, eventType string, equipment *dto.Equipment, log echo.Logger) {
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %v", err)
	}

	defer func() {
		if err := ch.Close(); err != nil {
			log.Errorf("Failed to close channel: %v", err)
		}
	}()

	body, err := json.Marshal(equipmentEvent{
		EventType: eventType,
		Equipment: *equipment,
	})
	if err != nil {
		log.Errorf("Failed to marshal event: %v", err)
	}

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

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Errorf("Failed to publish a message: %v", err)
	}
}
