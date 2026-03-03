package dto

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ==========================================
// Response Envelope Structures
// ==========================================

// APIResponse es el envelope estándar para todas las respuestas de la API
// @Description Estructura de respuesta estándar de la API
type APIResponse struct {
	Success   bool        `json:"success" example:"true"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp" example:"2024-01-15T10:30:00Z"`
} // @name APIResponse

// APIError representa un error en la respuesta
// @Description Detalles del error cuando success=false
type APIError struct {
	Code    string      `json:"code" example:"VALIDATION_ERROR"`
	Message string      `json:"message" example:"Los datos proporcionados no son válidos"`
	Details interface{} `json:"details,omitempty"`
} // @name APIError

// Meta contiene metadatos adicionales de la respuesta
// @Description Metadatos de paginación y otros
type Meta struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Count      int             `json:"count,omitempty"`
} // @name Meta

// ==========================================
// Legacy Structures (for backward compatibility)
// ==========================================

// ErrorResponse representa la respuesta estándar de error (legacy)
// @Description Respuesta de error estándar del API
type ErrorResponse struct {
	Error   string            `json:"error" example:"VALIDATION_ERROR"`
	Message string            `json:"message" example:"Validation failed"`
	Code    string            `json:"code,omitempty" example:"ERR_001"`
	Details map[string]string `json:"details,omitempty"`
} // @name ErrorResponse

// ValidationErrorResponse respuesta de error de validación con detalles
// @Description Respuesta de error cuando la validación falla
type ValidationErrorResponse struct {
	Error   string            `json:"error" example:"VALIDATION_ERROR"`
	Message string            `json:"message" example:"Validation failed"`
	Details map[string]string `json:"details" example:"email:Invalid email format,password:Password too short"`
} // @name ValidationErrorResponse

// SuccessResponse respuesta genérica exitosa (legacy)
// @Description Respuesta exitosa genérica
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
} // @name SuccessResponse

// PaginationMeta metadata de paginación
// @Description Información de paginación
type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
} // @name PaginationMeta

// PaginatedResponse respuesta paginada genérica
// @Description Respuesta con paginación
type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
} // @name PaginatedResponse

// HealthCheckResponse respuesta del endpoint de health check
// @Description Estado de salud del servicio
type HealthCheckResponse struct {
	Status    string `json:"status" example:"ok"`
	Database  string `json:"database" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2026-03-02T15:04:05Z"`
} // @name HealthCheckResponse

// ==========================================
// Helper Functions for Enterprise Envelope
// ==========================================

// Success crea una respuesta exitosa estándar con el nuevo envelope
func Success(c fiber.Ctx, status int, data interface{}) error {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(status).JSON(response)
}

// SuccessWithMeta crea una respuesta exitosa con metadatos
func SuccessWithMeta(c fiber.Ctx, status int, data interface{}, meta *Meta) error {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(status).JSON(response)
}

// Error crea una respuesta de error estándar con el nuevo envelope
func Error(c fiber.Ctx, status int, code string, message string, details interface{}) error {
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

// BadRequest respuesta para errores 400
func BadRequest(c fiber.Ctx, message string, details interface{}) error {
	return Error(c, fiber.StatusBadRequest, "BAD_REQUEST", message, details)
}

// Unauthorized respuesta para errores 401
func Unauthorized(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Forbidden respuesta para errores 403
func Forbidden(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFound respuesta para errores 404
func NotFound(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", message, nil)
}

// Conflict respuesta para errores 409
func Conflict(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, "CONFLICT", message, nil)
}

// ValidationError respuesta para errores de validación
func ValidationError(c fiber.Ctx, details interface{}) error {
	return Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Los datos proporcionados no son válidos", details)
}

// InternalServerError respuesta para errores 500
func InternalServerError(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Ha ocurrido un error interno del servidor"
	}
	return Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}

// ServiceUnavailable respuesta para errores 503
func ServiceUnavailable(c fiber.Ctx, message string) error {
	if message == "" {
		message = "El servicio no está disponible temporalmente"
	}
	return Error(c, fiber.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message, nil)
}

// TooManyRequests respuesta para errores de rate limiting
func TooManyRequests(c fiber.Ctx, retryAfter int) error {
	if retryAfter > 0 {
		c.Set("Retry-After", fmt.Sprintf("%d", retryAfter))
	}
	return Error(
		c,
		fiber.StatusTooManyRequests,
		"RATE_LIMIT_EXCEEDED",
		"Has excedido el límite de solicitudes. Intenta nuevamente más tarde.",
		nil,
	)
}
