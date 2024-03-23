package api

import (
	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/gin-gonic/gin"
)

func getAuthHandler(config conf.Config) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.APIUsername: config.APIPassword,
	})
}
