package api

import (
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
	Close()
}

func Run() error {
	conf, err := conf.LoadConfig()
	if err != nil {
		return types.WrapErr(err, "failed to load config")
	}

	var firestore FirestoreService
	firestore, err = repo.NewFirestoreService(conf)
	if err != nil {
		return types.WrapErr(err, "failed to create firestore service")
	}

	r := setupRouter(conf, firestore)

	return r.Run()
}

func setupRouter(conf conf.Config, fs FirestoreService) *gin.Engine {
	r := gin.Default()
	r.Use(headerMiddleware(conf))
	r.Use(requestIDMiddleware())

	// Auth endpoint, required to fetch a JWT
	r.POST("/auth", authHandler(conf))

	// Standard endpoints
	// All standard endpoints require a valid JWT
	authorized := r.Group("/", validateJWTMiddleware(conf))
	authorized.GET("/likes", getLikesHandler(fs))
	authorized.GET("/likes/:id", getLikeHandler(fs))
	authorized.GET("/posts", getPostsHandler(fs))

	return r
}
