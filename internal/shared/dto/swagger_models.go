package dto

// Este archivo contiene modelos adicionales para documentación Swagger

// QueryPaginationParams parámetros de paginación en query string
type QueryPaginationParams struct {
	// Número de página (default: 1)
	Page int `query:"page" example:"1" minimum:"1"`
	// Límite de registros por página (default: 10, max: 100)
	Limit int `query:"limit" example:"10" minimum:"1" maximum:"100"`
	// Campo por el cual ordenar
	Sort string `query:"sort" example:"created_at"`
	// Dirección de ordenamiento (asc o desc)
	Order string `query:"order" example:"desc" enums:"asc,desc"`
	// Término de búsqueda
	Search string `query:"search" example:"laptop"`
} // @name QueryPaginationParams

// QueryFilterParams parámetros de filtrado dinámico
type QueryFilterParams struct {
	// Filtro por estado: filter[status]=active
	FilterStatus string `query:"filter[status]" example:"active"`
	// Filtro mayor o igual: filter[price__gte]=100
	FilterPriceGte string `query:"filter[price__gte]" example:"100"`
	// Filtro menor o igual: filter[price__lte]=500
	FilterPriceLte string `query:"filter[price__lte]" example:"500"`
	// Filtro IN: filter[status__in]=active,pending
	FilterStatusIn string `query:"filter[status__in]" example:"active,pending"`
} // @name QueryFilterParams

// JWTClaims estructura de claims del JWT (solo para documentación)
type JWTClaims struct {
	// ID del usuario
	UserID string `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Email del usuario
	Email string `json:"email" example:"user@example.com"`
	// Tipo de usuario
	UserType string `json:"user_type" example:"company"`
	// Timestamp de expiración
	ExpiresAt int64 `json:"exp" example:"1709395200"`
	// Timestamp de emisión
	IssuedAt int64 `json:"iat" example:"1709308800"`
} // @name JWTClaims
