package middleware

import (
	"time"

	"github.com/atdevten/peace/internal/infrastructure/logging"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type LoggingMiddleware struct {
	logger *logging.Logger
}

func NewLoggingMiddleware(logger *logging.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get trace context
		ctx := c.Request.Context()
		span := trace.SpanFromContext(ctx)

		// Log request start
		logger := m.logger.WithFields(ctx, map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})
		logger.Info().Msg("Request started")

		// Add attributes to span
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request completion
		logger = m.logger.WithFields(ctx, map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": duration.String(),
			"size":     c.Writer.Size(),
		})
		logger.Info().Msg("Request completed")

		// Add response attributes to span
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int64("http.response_size", int64(c.Writer.Size())),
			attribute.String("http.duration", duration.String()),
		)
	}
}
