package main

import (
	"assign/core"
	"assign/handler"
	"assign/logger"
)

func Assign() {
	core.MakeDB()

	go core.Monitor()

	logger.ClearDB()
	handler.StartHandler()
	select {}
}
