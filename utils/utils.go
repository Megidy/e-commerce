package utils

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, message string, statusCode int) {
	if err != nil {
		c.JSON(statusCode, gin.H{
			"error":   err.Error(),
			"details": message,
		})
		log.Println(err)
		log.Println(err.Error())
	} else {
		c.JSON(statusCode, gin.H{
			"details": message,
		})
	}
}
func SendResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

func IsAccessory(id string) bool {
	return strings.HasPrefix(id, "a")
}
