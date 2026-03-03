package middleware

import (
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// RequestID middleware agrega un ID único a cada request
func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)

		return c.Next()
	}
}

// Logger middleware registra información de cada request
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Continuar con el siguiente middleware
		err := c.Next()

		// Calcular duración
		duration := time.Since(start)

		// Obtener información del request
		requestID, _ := c.Locals("requestID").(string)
		statusCode := c.Response().StatusCode()

		// Log del request
		logArgs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"status", statusCode,
			"duration_ms", duration.Milliseconds(),
			"request_id", requestID,
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		}

		if err != nil {
			logArgs = append(logArgs, "error", err.Error())
		}

		if statusCode >= 500 {
			logger.Error("Request completed with server error", logArgs...)
		} else if statusCode >= 400 {
			logger.Warn("Request completed with client error", logArgs...)
		} else {
			logger.Info("Request completed", logArgs...)
		}

		return err
	}
}

// CORS middleware configura Cross-Origin Resource Sharing
func CORS() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "3600")

		// Handle preflight request
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

// Security middleware agrega headers de seguridad
func Security() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		return c.Next()
	}
}
