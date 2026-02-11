package middleware

import (
	"english-learning/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware is a custom Gin middleware that logs request details using zap
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the log fields
		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if len(c.Errors) > 0 {
			// Log errors if any
			for _, e := range c.Errors.Errors() {
				logger.Log.Named("middleware").Error(e, fields...)
			}
		} else {
			// Log success/redirection/client error
			if status >= 500 {
				logger.Log.Named("middleware").Error("Internal Server Error", fields...)
			} else if status >= 400 {
				logger.Log.Named("middleware").Warn("Client Error", fields...)
			} else {
				logger.Log.Named("middleware").Info("Request Processed", fields...)
			}
		}
	}
}
