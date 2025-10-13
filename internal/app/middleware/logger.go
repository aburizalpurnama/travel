package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger adalah middleware untuk membuat contextual logger dan mencatat ringkasan request.
func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqID := uuid.New().String()
		reqLogger := logger.With("request_id", reqID)
		c.Locals("logger", reqLogger)

		err := c.Next()

		latency := time.Since(start)
		statusCode := c.Response().StatusCode()
		clientIP := c.IP()

		requestError := c.Locals("error")

		if requestError != nil {
			if e, ok := requestError.(error); ok {
				reqLogger.Error(
					"Request failed with an error",
					"method", c.Method(),
					"path", c.Path(),
					"status_code", statusCode,
					"latency_ms", latency.Milliseconds(),
					"client_ip", clientIP,
					"error", e.Error(),
				)
			}
		} else {
			if statusCode >= 500 {
				reqLogger.Error(
					"Request completed with server error",
					"method", c.Method(), "path", c.Path(), "status_code", statusCode,
					"latency_ms", latency.Milliseconds(), "client_ip", clientIP,
				)
			} else if statusCode >= 400 {
				reqLogger.Warn(
					"Request completed with client error",
					"method", c.Method(), "path", c.Path(), "status_code", statusCode,
					"latency_ms", latency.Milliseconds(), "client_ip", clientIP,
				)
			} else {
				reqLogger.Info(
					"Request handled successfully",
					"method", c.Method(), "path", c.Path(), "status_code", statusCode,
					"latency_ms", latency.Milliseconds(), "client_ip", clientIP,
				)
			}
		}

		return err
	}
}
