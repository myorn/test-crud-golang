package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"equipment-registry-service/controllers"
	"equipment-registry-service/db"
	"equipment-registry-service/mq"
)

func main() {
	// init DB and migrate schema
	controller := &controllers.Controller{
		DB: db.Init(),
		MQ: mq.Connect(),
	}

	defer controller.MQ.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/equipment", controller.CreateEquipment)
	e.GET("/equipment/:id", controller.GetEquipment)
	e.PUT("/equipment/:id", controller.UpdateEquipment)
	e.DELETE("/equipment/:id", controller.DeleteEquipment)
	e.PATCH("/equipment/restore/:id", controller.RestoreEquipment)

	e.Logger.Fatal(e.Start(":8080"))
}
