package errors

import (
	"errors"
	"fmt"
)

// ErrorType representa el tipo de error
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeConflict       ErrorType = "CONFLICT"
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeDatabase       ErrorType = "DATABASE_ERROR"
	ErrorTypeAuthentication ErrorType = "AUTHENTICATION_ERROR"
)

// AppError representa un error de aplicación personalizado
type AppError struct {
	Type       ErrorType         `json:"type"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	StatusCode int               `json:"-"`
	Err        error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError crea un nuevo error de aplicación
func NewAppError(errType ErrorType, message string, statusCode int) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithError agrega un error subyacente
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithDetails agrega detalles adicionales al error
func (e *AppError) WithDetails(details map[string]string) *AppError {
	e.Details = details
	return e
}

// Errores predefinidos comunes
var (
	ErrValidation = NewAppError(
		ErrorTypeValidation,
		"Validation failed",
		400,
	)

	ErrNotFound = NewAppError(
		ErrorTypeNotFound,
		"Resource not found",
		404,
	)

	ErrUnauthorized = NewAppError(
		ErrorTypeUnauthorized,
		"Unauthorized access",
		401,
	)

	ErrForbidden = NewAppError(
		ErrorTypeForbidden,
		"Forbidden",
		403,
	)

	ErrConflict = NewAppError(
		ErrorTypeConflict,
		"Resource already exists",
		409,
	)

	ErrInternal = NewAppError(
		ErrorTypeInternal,
		"Internal server error",
		500,
	)

	ErrBadRequest = NewAppError(
		ErrorTypeBadRequest,
		"Bad request",
		400,
	)

	ErrDatabase = NewAppError(
		ErrorTypeDatabase,
		"Database error",
		500,
	)

	ErrAuthentication = NewAppError(
		ErrorTypeAuthentication,
		"Authentication failed",
		401,
	)
)

// IsAppError verifica si un error es de tipo AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extrae el AppError de un error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
