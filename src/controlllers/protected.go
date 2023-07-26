package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProtectedController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "You need to be authenticated to see this.",
	})
}
