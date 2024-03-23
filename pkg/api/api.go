package api

import (
	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/repo"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
)

func Run() error {
	config, err := conf.LoadConfig()
	if err != nil {
		return types.WrapErr(err, "failed to load config")
	}

	_, err = repo.NewFirestoreService(config)
	if err != nil {
		return types.WrapErr(err, "failed to create firestore service")
	}

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello",
		})
	})

	return r.Run()
}
