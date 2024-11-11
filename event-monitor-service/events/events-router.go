package events

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"event-monitor-service/storage"
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

	Equipment
}

// Equipment represents the equipment entity
type Equipment struct {
	ID         uuid.UUID      `json:"id"`
	Type       string         `json:"type"`
	Status     string         `json:"status"`
	Parameters map[string]any `json:"parameters"`
}

func Router(body []byte) {
	event := new(equipmentEvent)

	if err := json.Unmarshal(body, event); err != nil {
		log.Errorf("Failed to unmarshal message: %v", err)

		return
	}

	switch event.EventType {
	case EquipmentCreated:
		TriggerOnboarding()

	case EquipmentRead:

	case EquipmentUpdated:
		RefreshConfiguration()
		UpdateSubsystems()

	case EquipmentDeleted:
		Shutdown()

	case EquipmentRestored:
		Reinstall()

	default:
		log.Errorf("Unknown event type: %s", event.EventType)
	}

	storage.StoreEvent(event.EventType, event.Equipment)
}
