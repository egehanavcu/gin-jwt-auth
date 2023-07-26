package middlewares

import (
	"gin-jwt-auth/src/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct{}

func NewJWTMiddleware() *JWTMiddleware {
	return &JWTMiddleware{}
}

func (m *JWTMiddleware) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		token := parts[1]
		data, err := services.ParseJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// Check if token has expired
		if time.Now().Unix() > data["exp"].(int64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token has expired",
			})
			return
		}

		c.Set("user", data)
		c.Next()
	}
}
