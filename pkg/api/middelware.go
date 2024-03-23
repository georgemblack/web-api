package api

import (
	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/gin-gonic/gin"
)

func getHeaderMiddleware(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.AllowedOriginHeader)
		c.Header("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	}
}
