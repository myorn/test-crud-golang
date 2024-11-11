package main

import (
	"event-monitor-service/mq"
	"event-monitor-service/storage"
)

func main() {
	conn := mq.Connect()

	storage.InitEventLog()
	mq.ReceiverLoop(conn)
}
