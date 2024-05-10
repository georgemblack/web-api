package api

import (
	"net/http"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/repo"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
)

type FirestoreService interface {
	GetLike(id string) (types.Like, error)
	GetLikes() ([]types.Like, error)
	AddLike(like types.Like) (string, error)
	DeleteLike(id string) error
	GetPost(id string) (types.Post, error)
	GetPosts(filters repo.PostFilters) ([]types.Post, error)
	AddPost(post types.Post) (string, error)
	UpdatePost(post types.Post) error
	DeletePost(id string) error
	GetHashList() (types.HashList, error)
	UpdateHashList(hashList types.HashList) error
	Close()
}

func Run() error {
	config, err := conf.LoadConfig()
	if err != nil {
		return types.WrapErr(err, "failed to load config")
	}

	var firestore FirestoreService
	firestore, err = repo.NewFirestoreService(config)
	if err != nil {
		return types.WrapErr(err, "failed to create firestore service")
	}

	r := setupRouter(config, firestore)

	return r.Run()
}

func setupRouter(config conf.Config, firestore FirestoreService) *gin.Engine {
	r := gin.Default()
	r.Use(headerMiddleware(config))
	r.Use(requestIDMiddleware())

	// Auth endpoint, required to fetch a JWT
	r.POST("/auth", authHandler(config))

	// Standard endpoints
	// All standard endpoints require a valid JWT
	authorized := r.Group("/", validateJWTMiddleware(config))

	authorized.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello",
		})
	})

	authorized.GET("/likes", func(c *gin.Context) {
		likes, err := firestore.GetLikes()
		if err != nil {
			internalServerError(c)
		}
		c.JSON(http.StatusOK, likes)
	})

	return r
}
