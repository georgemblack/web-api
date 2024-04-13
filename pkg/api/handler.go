package api

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getAuthHandler(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokens := strings.Split(header, " ")
		if len(tokens) != 2 {
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(tokens[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			return
		}
		credentials := strings.Split(string(decoded), ":")
		if len(credentials) != 2 {
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			return
		}
		if credentials[0] != config.APIUsername || credentials[1] != config.APIPassword {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Request is valid, return signed JWT
		expires := time.Now().Add(time.Hour * 6)
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtTokenStr, err := jwtToken.SignedString([]byte(config.TokenSecret))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to sign token"})
			return
		}
		c.JSON(200, gin.H{"token": jwtTokenStr})
	}
}
