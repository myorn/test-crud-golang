package storage

import (
	"os"

	"github.com/labstack/gommon/log"
)

var fileLogger *log.Logger

func InitEventLog() {
	const permission = 0666

	file, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, permission)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	fileLogger = log.New("fileLogger")
	fileLogger.SetOutput(file)
	fileLogger.SetLevel(log.INFO)
	fileLogger.SetHeader("${time_rfc3339}")
}

func StoreEvent(eventType string, payload interface{}) {
	fileLogger.Infof("%s: %+v", eventType, payload)
}
