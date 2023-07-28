package controllers

import (
	"gin-jwt-auth/src/config"
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
		userData := map[string]interface{}{
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(config.Get("AccessTokenDuration").(time.Duration)).Unix(),
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

		refreshData := map[string]interface{}{
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(config.Get("RefreshTokenDuration").(time.Duration)).Unix(),
			config.Get("AccessTokenIdentifierKey").(string): userData[config.Get("AccessTokenIdentifierKey").(string)],
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
			"message":      "Login successful",
			"access_token": access_token,
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

func RefreshHandler(c *gin.Context) {
	access_token, access_err := c.Cookie("access_token")
	refresh_token, refresh_err := c.Cookie("refresh_token")

	if access_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access token not found",
		})
		return
	}

	if refresh_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Refresh token not found",
		})
		return
	}

	access_data, access_data_err := services.HandleJWT(access_token)
	if access_data_err != nil && access_data_err.Error() != "token has expired" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": access_data_err.Error(),
		})
		return
	}

	if access_data_err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access token has not expired yet",
		})
		return
	}

	new_access_data := make(map[string]interface{})
	for k, v := range access_data.(map[string]interface{}) {
		new_access_data[k] = v
	}
	new_access_data["iat"] = time.Now().Unix()
	new_access_data["exp"] = time.Now().Add(config.Get("AccessTokenDuration").(time.Duration)).Unix()

	refresh_data, refresh_data_err := services.HandleJWT(refresh_token)

	if refresh_data_err != nil && refresh_data_err.Error() == "token has expired" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Refresh token has expired",
		})
		return
	}

	if refresh_data.(map[string]interface{})[config.Get("AccessTokenIdentifierKey").(string)] != access_data.(map[string]interface{})[config.Get("AccessTokenIdentifierKey").(string)] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Refresh token is not valid",
		})
		return
	}

	new_access_token, err := services.GenerateJWT(new_access_data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", new_access_token, 3600, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Token refreshed",
		"access_token": new_access_data,
	})
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, false)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
