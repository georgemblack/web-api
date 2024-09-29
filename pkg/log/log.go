package log

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Warn(c *gin.Context, message string) {
	slog.Warn(message, "requestId", c.GetString("requestId"))
}

func Error(c *gin.Context, message string) {
	slog.Error(message, "requestId", c.GetString("requestId"))
}
