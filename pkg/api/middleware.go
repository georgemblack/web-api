package api

import (
	"strings"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func headerMiddleware(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.AllowedOriginHeader)
		c.Header("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	}
}

func validateJWTMiddleware(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token from header
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")
		if len(parts) != 2 {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		token := parts[1]

		// Parse & validate token
		_, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return []byte(config.TokenSecret), nil
		})
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
