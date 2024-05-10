package types

import "github.com/gin-gonic/gin"

type GlobalConfig struct {
	CRI string
	Gin *gin.Engine
}
