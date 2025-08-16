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

		level := logger.GetLevel()

		if len(c.Errors) > 0 && level <= zerolog.ErrorLevel {
			for _, err := range c.Errors {
				event := logger.Error()

				event.Err(err.Err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Int("status", c.Writer.Status()).
					Dur("duration", duration).
					Str("ip", c.ClientIP()).
					Str("user_agent", c.Request.UserAgent())

				if logger.GetLevel() <= zerolog.TraceLevel {
					event.Str("trace", fmt.Sprintf("%+v", err.Err))
				}

				event.Msg("HTTP request error")
			}
		} else if c.Writer.Status() >= 400 && level <= zerolog.WarnLevel {
			event := logger.Warn()

			event.Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("duration", duration).
				Str("ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent())

			event.Msg("HTTP request completed with warning")
		} else if level <= zerolog.InfoLevel {
			event := logger.Info()

			event.Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("duration", duration).
				Str("ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent())

			event.Msg("HTTP request completed")
		}
	}
}
