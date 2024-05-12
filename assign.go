package main

import (
	"assign/core"
	"assign/handler"
)

func Assign() {
	core.MakeDB()

	// 로그인 접속(인증, ssh 접속?)(프론트)

	// 이후
	// container, service 관리(생성)
	// container, service 관리(추적, 삭제)
	handler.StartHandler()
	for {

	}
}
