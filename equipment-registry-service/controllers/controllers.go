package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	dao "equipment-registry-service/db"
	"equipment-registry-service/dto"
	"equipment-registry-service/mq"
)

type Controller struct {
	DB *gorm.DB
	MQ *amqp.Connection
}

func (ctrl *Controller) CreateEquipment(c echo.Context) error {
	ctx := c.Request().Context()

	equipment := new(dto.Equipment)

	if err := c.Bind(equipment); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	equipment.ID = uuid.New()

	if err := dao.EquipmentCreateOrChange(ctx, ctrl.DB, &dao.EquipmentTable{
		ID:         equipment.ID,
		Type:       equipment.Type,
		Status:     equipment.Status,
		Parameters: equipment.Parameters,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// send event
	mq.SendEvent(ctx, ctrl.MQ, mq.EquipmentCreated, equipment, c.Logger())

	return c.JSON(http.StatusCreated, equipment)
}

func (ctrl *Controller) GetEquipment(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	equipmentDB, err := dao.EquipmentByID(ctx, ctrl.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Equipment not found")
		}
	}

	equipment := &dto.Equipment{
		ID:         equipmentDB.ID,
		Type:       equipmentDB.Type,
		Status:     equipmentDB.Status,
		Parameters: equipmentDB.Parameters,
	}

	// send event
	mq.SendEvent(ctx, ctrl.MQ, mq.EquipmentRead, equipment, c.Logger())

	return c.JSON(http.StatusOK, equipment)
}

func (ctrl *Controller) UpdateEquipment(c echo.Context) error {
	ctx := c.Request().Context()

	equipment := new(dto.Equipment)

	if err := c.Bind(equipment); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// let's pretend we have a high load
	err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
		// check if equipment exists
		_, err := dao.EquipmentByIDWithLock(ctx, tx, equipment.ID)
		if err != nil {
			return err
		}

		// update equipment
		if err := dao.EquipmentCreateOrChange(ctx, tx, &dao.EquipmentTable{
			ID:         equipment.ID,
			Type:       equipment.Type,
			Status:     equipment.Status,
			Parameters: equipment.Parameters,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Equipment not found")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	// send event
	mq.SendEvent(ctx, ctrl.MQ, mq.EquipmentUpdated, equipment, c.Logger())

	return c.JSON(http.StatusOK, equipment)
}

func (ctrl *Controller) DeleteEquipment(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	if err := dao.EquipmentDeleteByID(ctx, ctrl.DB, id); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	equipmentDB, err := dao.EquipmentByIDUnscoped(ctx, ctrl.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Equipment not found")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	// send event
	mq.SendEvent(ctx, ctrl.MQ, mq.EquipmentDeleted, &dto.Equipment{
		ID:         equipmentDB.ID,
		Type:       equipmentDB.Type,
		Status:     equipmentDB.Status,
		Parameters: equipmentDB.Parameters,
	}, c.Logger())

	return c.NoContent(http.StatusNoContent)
}

func (ctrl *Controller) RestoreEquipment(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	if err := dao.EquipmentRestoreByID(ctx, ctrl.DB, id); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	equipmentDB, err := dao.EquipmentByID(ctx, ctrl.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Equipment not found")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	// send event
	mq.SendEvent(ctx, ctrl.MQ, mq.EquipmentRestored, &dto.Equipment{
		ID:         equipmentDB.ID,
		Type:       equipmentDB.Type,
		Status:     equipmentDB.Status,
		Parameters: equipmentDB.Parameters,
	}, c.Logger())

	return c.NoContent(http.StatusNoContent)
}
