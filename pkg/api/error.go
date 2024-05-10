package api

import (
	"net/http"
	"time"

	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
)

func internalServerError(c *gin.Context) {
	resp := types.ErrorResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "Internal server error",
		RequestID: c.GetString("requestId"),
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
}
