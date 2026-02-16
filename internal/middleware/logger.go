package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)

		logger.Info().
			Str("ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Msg("request")
	}
}
