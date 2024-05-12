package handler

import (
	"assign/common"
	"assign/config"
	"assign/core"
	"assign/logger"
	"net/http"

	"assign/types"

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

func StartHandler() error {
	config.GlobalConfig.Gin.Use(CORSMiddleware())

	config.GlobalConfig.Gin.GET("/login", login)
	manager := config.GlobalConfig.Gin.Group("/Manager")
	{
		regist := manager.Group("/regist")
		{
			regist.PUT("/person", registPerson)
		}
		delete := manager.Group("/delete")
		{
			delete.DELETE("/person", deletePerson)
		}
	}

	config.GlobalConfig.Gin.Run(":5000")
	return nil
}

/////////////
/*  func   */
/////////////

func login(c *gin.Context) {
	var person types.Person
	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if person.ID == "" || person.Position == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID and Position are required"})
		return
	}

	if person.Position == "manager" {
		if person.ID == "nsm" {
			c.JSON(http.StatusOK, gin.H{"message": "Manager login successfully"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	} else if person.Position == "user" {
		if ok := logger.IsExist(person.ID); !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			info := logger.GetInfo(person.ID)
			c.JSON(http.StatusOK, gin.H{"info": info})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position"})
	}
}

func registPerson(c *gin.Context) {
	var person types.Person
	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if person.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := logger.RegistPerson(person.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register person"})
		return
	} else {
		info := core.MakeServer()
		info.StudentID = person.ID
		logger.UpdateInfo(info)
		info.Ip = common.GetServerIP()
		c.JSON(http.StatusOK, gin.H{"message": "Person registered successfully", "server_info": info})
	}
}

func deletePerson(c *gin.Context) {
	var person types.Person
	if person.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	err := logger.DeletePerson(person.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}
