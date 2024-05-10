package core

import "assign/config"

func Assign() {
	CreateServiceAndContainer(config.GlobalConfig.KubernetesClient)
	// db container 생성
	// db container 식별
	// container 관리(생성, db 관리)
	// container 관리(추적, 삭제, db 관리)
	// 로그인 접속(인증, ssh 접속?)
}
