package api

import (
	"strings"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/log"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Populates a request context with a randomly generated ID that can be referenced throughout the lifetime of the request.
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requestId", uuid.New().String())
	}
}

// Populates a request with required CORS headers.
func headerMiddleware(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.AllowedOriginHeader)
		c.Header("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	}
}

// Validates a JWT in the request header.
func validateJWTMiddleware(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token from header
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")
		if len(parts) != 2 {
			log.Warn(c, "faild to parse two parts of 'Authorization' header")
			unauthorizedError(c)
			return
		}
		token := parts[1]

		// Parse & validate token
		_, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return []byte(config.TokenSecret), nil
		})
		if err != nil {
			log.Warn(c, types.WrapErr(err, "invalid jwt").Error())
			unauthorizedError(c)
			return
		}

		c.Next()
	}
}
