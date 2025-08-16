package mervermid

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GinZeroLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error().
					Err(err.Err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Str("trace", fmt.Sprintf("%+v", err.Err)).
					Int("status", c.Writer.Status()).
					Dur("duration", duration).
					Msg("HTTP request error")
			}
		} else if c.Writer.Status() >= 400 {
			logger.Warn().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("duration", duration).
				Msg("HTTP request completed with warning")
		} else {
			logger.Info().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("duration", duration).
				Msg("HTTP request completed")
		}
	}
}
