package dto

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// APIResponse es el envelope estándar para todas las respuestas de la API
// @Description Estructura de respuesta estándar de la API
type APIResponse struct {
	Success   bool        `json:"success" example:"true"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp" example:"2024-01-15T10:30:00Z"`
}

// APIError representa un error en la respuesta
// @Description Detalles del error cuando success=false
type APIError struct {
	Code    string      `json:"code" example:"VALIDATION_ERROR"`
	Message string      `json:"message" example:"Los datos proporcionados no son válidos"`
	Details interface{} `json:"details,omitempty"`
}

// Meta contiene metadatos adicionales de la respuesta
// @Description Metadatos de paginación y otros
type Meta struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Count      int             `json:"count,omitempty"`
}

// SuccessResponse crea una respuesta exitosa estándar
func SuccessResponse(c fiber.Ctx, status int, data interface{}) error {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(status).JSON(response)
}

// SuccessResponseWithMeta crea una respuesta exitosa con metadatos
func SuccessResponseWithMeta(c fiber.Ctx, status int, data interface{}, meta *Meta) error {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(status).JSON(response)
}

// ErrorResponse crea una respuesta de error estándar
func ErrorResponse(c fiber.Ctx, status int, code string, message string, details interface{}) error {
	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(status).JSON(response)
}

// BadRequestError respuesta para errores 400
func BadRequestError(c fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "BAD_REQUEST", message, details)
}

// UnauthorizedError respuesta para errores 401
func UnauthorizedError(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// ForbiddenError respuesta para errores 403
func ForbiddenError(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFoundError respuesta para errores 404
func NotFoundError(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, "NOT_FOUND", message, nil)
}

// ConflictError respuesta para errores 409
func ConflictError(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, "CONFLICT", message, nil)
}

// ValidationError respuesta para errores de validación
func ValidationError(c fiber.Ctx, details interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Los datos proporcionados no son válidos", details)
}

// InternalServerError respuesta para errores 500
func InternalServerError(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Ha ocurrido un error interno del servidor"
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}

// ServiceUnavailableError respuesta para errores 503
func ServiceUnavailableError(c fiber.Ctx, message string) error {
	if message == "" {
		message = "El servicio no está disponible temporalmente"
	}
	return ErrorResponse(c, fiber.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message, nil)
}

// TooManyRequestsError respuesta para errores de rate limiting
func TooManyRequestsError(c fiber.Ctx, retryAfter int) error {
	if retryAfter > 0 {
		c.Set("Retry-After", string(rune(retryAfter)))
	}
	return ErrorResponse(
		c,
		fiber.StatusTooManyRequests,
		"RATE_LIMIT_EXCEEDED",
		"Has excedido el límite de solicitudes. Intenta nuevamente más tarde.",
		nil,
	)
}
