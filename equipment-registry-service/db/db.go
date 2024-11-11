package db

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// EquipmentTable represents the equipment entity
type EquipmentTable struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Type       string
	Status     string
	Parameters JSONB          `gorm:"type:jsonb"`
	DeletedAt  gorm.DeletedAt `gorm:"column:DELETED_AT;type:timestamp"`
}

func Init() *gorm.DB {
	dsn := "host=postgres2 user=postgres dbname=postgres sslmode=disable password=postgres"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// migrate the schema
	if err := db.AutoMigrate(&EquipmentTable{}); err != nil {
		log.Fatalf("failed to migrate table: %v", err)
	}

	log.Println("Successfully connected to database and migrated schema")

	return db
}
