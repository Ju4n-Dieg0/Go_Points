package dto

// ErrorResponse representa la respuesta estándar de error
// @Description Respuesta de error estándar del API
type ErrorResponse struct {
	// Tipo de error
	Error string `json:"error" example:"VALIDATION_ERROR"`
	// Mensaje descriptivo del error
	Message string `json:"message" example:"Validation failed"`
	// Código de error específico (opcional)
	Code string `json:"code,omitempty" example:"ERR_001"`
	// Detalles adicionales del error (opcional)
	Details map[string]string `json:"details,omitempty"`
} // @name ErrorResponse

// ValidationErrorResponse respuesta de error de validación con detalles
// @Description Respuesta de error cuando la validación falla
type ValidationErrorResponse struct {
	// Tipo de error
	Error string `json:"error" example:"VALIDATION_ERROR"`
	// Mensaje general
	Message string `json:"message" example:"Validation failed"`
	// Detalles de cada campo con error
	Details map[string]string `json:"details" example:"email:Invalid email format,password:Password too short"`
} // @name ValidationErrorResponse

// SuccessResponse respuesta genérica exitosa
// @Description Respuesta exitosa genérica
type SuccessResponse struct {
	// Mensaje de éxito
	Message string `json:"message" example:"Operation completed successfully"`
	// Datos adicionales (opcional)
	Data interface{} `json:"data,omitempty"`
} // @name SuccessResponse

// PaginationMeta metadata de paginación
// @Description Información de paginación
type PaginationMeta struct {
	// Página actual
	Page int `json:"page" example:"1"`
	// Límite de registros por página
	Limit int `json:"limit" example:"10"`
	// Total de registros
	Total int64 `json:"total" example:"100"`
	// Total de páginas
	TotalPages int `json:"total_pages" example:"10"`
} // @name PaginationMeta

// PaginatedResponse respuesta paginada genérica
// @Description Respuesta con paginación
type PaginatedResponse struct {
	// Datos paginados
	Data interface{} `json:"data"`
	// Metadata de paginación
	Meta PaginationMeta `json:"meta"`
} // @name PaginatedResponse

// HealthCheckResponse respuesta del endpoint de health check
// @Description Estado de salud del servicio
type HealthCheckResponse struct {
	// Estado general
	Status string `json:"status" example:"ok"`
	// Estado de la base de datos
	Database string `json:"database" example:"healthy"`
	// Timestamp
	Timestamp string `json:"timestamp" example:"2026-03-02T15:04:05Z"`
} // @name HealthCheckResponse
