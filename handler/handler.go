package handler

import (
	"assign/config"

	"github.com/gin-gonic/gin"
)

////////////////
/*  Handler   */
////////////////

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// StartHandler starts the server with the configured routes.
func StartHandler() error {
	config.GlobalConfig.Gin.Use(CORSMiddleware())

	// block := config.GlobalConfig.Gin.Group("/block")
	// {
	// 	ip := block.Group("/ip")
	// 	{
	// 		ip.GET("/list", blockIPv4List)
	// 		ip.PUT("/on", blockIPv4BlockOn)
	// 		ip.PUT("/off", blockIPv4BlockOff)
	// 	}
	// 	mac := block.Group("/mac")
	// 	{
	// 		mac.GET("/list", blockMACList)
	// 		mac.PUT("/on", blockMACBlockOn)
	// 		mac.PUT("/off", blockMACBlockOff)
	// 	}
	// }

	config.GlobalConfig.Gin.Run(":8081")
	return nil
}

/////////////
/*  func   */
/////////////

func isExistingUser(c *gin.Context) {
}
