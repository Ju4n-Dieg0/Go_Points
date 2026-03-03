package errors

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"github.com/gofiber/fiber/v3"
)

// ErrorResponse representa la estructura de respuesta de error
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Type    string            `json:"type,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorHandler es el middleware global de manejo de errores para Fiber v3
func ErrorHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Continuar con la siguiente función
		err := c.Next()

		if err == nil {
			return nil
		}

		// Obtener el error original
		var appErr *AppError
		statusCode := fiber.StatusInternalServerError
		response := ErrorResponse{
			Success: false,
		}

		// Verificar si es un AppError personalizado
		if GetAppError(err) != nil {
			appErr = GetAppError(err)
			statusCode = appErr.StatusCode
			response.Error = appErr.Message
			response.Type = string(appErr.Type)
			response.Details = appErr.Details

			// Logging basado en severidad
			if statusCode >= 500 {
				logger.Error("Internal server error",
					"error", err.Error(),
					"type", appErr.Type,
					"path", c.Path(),
					"method", c.Method(),
				)
			} else {
				logger.Warn("Client error",
					"error", appErr.Message,
					"type", appErr.Type,
					"path", c.Path(),
					"method", c.Method(),
					"status", statusCode,
				)
			}
		} else {
			// Manejar errores de Fiber
			if fiberErr, ok := err.(*fiber.Error); ok {
				statusCode = fiberErr.Code
				response.Error = fiberErr.Message
				response.Type = "FIBER_ERROR"

				logger.Warn("Fiber error",
					"error", fiberErr.Message,
					"path", c.Path(),
					"method", c.Method(),
					"status", statusCode,
				)
			} else {
				// Error desconocido
				response.Error = "An unexpected error occurred"
				response.Type = string(ErrorTypeInternal)

				logger.Error("Unexpected error",
					"error", err.Error(),
					"path", c.Path(),
					"method", c.Method(),
				)
			}
		}

		// Enviar respuesta de error
		return c.Status(statusCode).JSON(response)
	}
}

// RecoverMiddleware captura panics y los convierte en errores manejables
func RecoverMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered",
					"panic", r,
					"path", c.Path(),
					"method", c.Method(),
				)

				_ = c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Success: false,
					Error:   "Internal server error",
					Type:    string(ErrorTypeInternal),
				})
			}
		}()

		return c.Next()
	}
}
