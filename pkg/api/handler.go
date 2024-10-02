package api

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/log"
	"github.com/georgemblack/web-api/pkg/repo"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authHandler(config conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokens := strings.Split(header, " ")
		if len(tokens) != 2 {
			log.Warn(c, "failed to parse two words in 'Authorization' header")
			unauthorizedError(c)
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(tokens[1])
		if err != nil {
			log.Warn(c, "failed to decode encoded username/password in 'Authorization' header")
			unauthorizedError(c)
			return
		}
		credentials := strings.Split(string(decoded), ":")
		if len(credentials) != 2 {
			log.Warn(c, "failed to parse colon-separated username/password in 'Authorization' header")
			unauthorizedError(c)
			return
		}
		if credentials[0] != config.APIUsername || credentials[1] != config.APIPassword {
			log.Warn(c, "invalid credentials provided")
			unauthorizedError(c)
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
			log.Error(c, "failed to sign jwt token")
			internalServerError(c)
			return
		}
		c.JSON(http.StatusOK, types.AuthResponse{Token: jwtTokenStr})
	}
}

func getLikesHandler(fs FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		likes, err := fs.GetLikes()
		if err != nil {
			log.Error(c, types.WrapErr(err, "failed to get likes").Error())
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
			log.Warn(c, "'id' param unexpectedly empty")
			invalidRequestError(c)
			return
		}

		like, err := fs.GetLike(id)
		if err != nil {
			log.Error(c, types.WrapErr(err, "failed to get like").Error())
			internalServerError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"like": like})
	}
}

func getPostsHandler(fs FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := repo.PostFilters{}
		published := c.Query("published")
		listed := c.Query("listed")

		t := true
		f := false
		if published == "true" {
			filters.Published = &t
		}
		if published == "false" {
			filters.Published = &f
		}
		if listed == "true" {
			filters.Listed = &t
		}
		if listed == "false" {
			filters.Listed = &f
		}

		posts, err := fs.GetPosts(filters)
		if err != nil {
			log.Error(c, types.WrapErr(err, "failed to get posts").Error())
			internalServerError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"posts": posts})
	}
}

func updatePostHandler(fs FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			log.Warn(c, "'id' param unexpectedly empty")
			invalidRequestError(c)
			return
		}

		var post types.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			log.Warn(c, types.WrapErr(err, "failed to bind json").Error())
			invalidRequestError(c)
			return
		}

		post.ID = id
		if err := fs.UpdatePost(post); err != nil {
			log.Error(c, types.WrapErr(err, "failed to update post").Error())
			internalServerError(c)
			return
		}
		c.Status(http.StatusOK)
	}
}
