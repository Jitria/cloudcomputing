package main

import (
	"assign/core"
	"assign/handler"
	"assign/logger"
)

func Assign() {
	core.Monitor()
	core.MakeDB()

	// 로그인 접속(인증, ssh 접속?)(프론트)

	// 이후
	// container, service 관리(생성)
	// container, service 관리(추적, 삭제)
	logger.ClearDB()
	handler.StartHandler()
	for {

	}
}
