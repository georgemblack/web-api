package api

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authHandler(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokens := strings.Split(header, " ")
		if len(tokens) != 2 {
			slog.Warn("failed to parse two words in 'Authorization' header")
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(tokens[1])
		if err != nil {
			slog.Warn("failed to decode encoded username/password in 'Authorization' header")
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			return
		}
		credentials := strings.Split(string(decoded), ":")
		if len(credentials) != 2 {
			slog.Warn("failed to parse colon-separated username/password in 'Authorization' header")
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			return
		}
		if credentials[0] != config.APIUsername || credentials[1] != config.APIPassword {
			slog.Warn("invalid credentials provided")
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Request is valid, return signed JWT
		expires := time.Now().Add(time.Hour * 2)
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtTokenStr, err := jwtToken.SignedString([]byte(config.TokenSecret))
		if err != nil {
			slog.Error("failed to sign jwt token")
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(200, types.AuthResponse{Token: jwtTokenStr})
	}
}

func getLikesHandler(fs FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		likes, err := fs.GetLikes()
		if err != nil {
			internalServerError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"likes": likes})
	}
}

func getLikeHandler(fs FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			invalidRequestError(c)
			return
		}

		like, err := fs.GetLike(id)
		if err != nil {
			internalServerError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"like": like})
	}
}
