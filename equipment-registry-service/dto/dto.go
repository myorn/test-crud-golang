package dto

import (
	"github.com/google/uuid"
)

// Equipment represents the equipment entity
type Equipment struct {
	ID         uuid.UUID      `json:"id"`
	Type       string         `json:"type"`
	Status     string         `json:"status"`
	Parameters map[string]any `json:"parameters"`
}
