package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RequireStaticToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Static-Token")
		expected := os.Getenv("STATIC_TOKEN_FOR_N8N")
		if token == "" || token != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing static token"})
			return
		}
		c.Next()
	}
}
