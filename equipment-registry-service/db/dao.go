package db

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func EquipmentCreateOrChange(ctx context.Context, tx *gorm.DB, equipment *EquipmentTable) error {
	return tx.WithContext(ctx).Save(equipment).Error
}

func EquipmentByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*EquipmentTable, error) {
	equipment := &EquipmentTable{}

	err := tx.
		WithContext(ctx).
		Where("ID = ?", id).
		Take(&equipment).
		Error

	if err != nil {
		return nil, err
	}

	return equipment, nil
}

func EquipmentByIDUnscoped(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*EquipmentTable, error) {
	equipment := &EquipmentTable{}

	err := tx.
		WithContext(ctx).
		Unscoped().
		Where("ID = ?", id).
		Take(&equipment).
		Error

	if err != nil {
		return nil, err
	}

	return equipment, nil
}

func EquipmentByIDWithLock(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*EquipmentTable, error) {
	equipment := &EquipmentTable{}

	// select for update
	err := tx.
		WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Take(equipment, id).
		Error

	if err != nil {
		return nil, err
	}

	return equipment, nil
}

func EquipmentDeleteByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	return tx.
		WithContext(ctx).
		Where("ID = ?", id).
		Delete(&EquipmentTable{}).
		Error
}

func EquipmentRestoreByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	return tx.
		WithContext(ctx).
		Unscoped().
		Model(&EquipmentTable{}).
		Where("ID = ?", id).
		Update("DELETED_AT", nil).
		Error
}
