package controllers

import (
	"gin-jwt-auth/src/dto"
	"gin-jwt-auth/src/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var loginDTO dto.AuthLoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if loginDTO.Email == "admin@admin.com" && loginDTO.Password == "admin" {
		userData := map[string]any{
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
			"email": loginDTO.Email,
			"role":  "admin",
		}

		access_token, err := services.GenerateJWT(userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		refreshData := map[string]any{
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"token": access_token,
		}

		refresh_token, err := services.GenerateJWT(refreshData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.SetCookie("access_token", access_token, 3600, "/", "", false, false)
		c.SetCookie("refresh_token", refresh_token, 3600, "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"access_token":  access_token,
			"refresh_token": refresh_token,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
	}
}

func RegisterHandler(c *gin.Context) {
	// TODO: Implement registration logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
	})
}

func LogoutHandler(c *gin.Context) {
	// TODO: Implement logout logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
