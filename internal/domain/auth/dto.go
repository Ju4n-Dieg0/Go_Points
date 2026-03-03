package auth

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest DTO para registro de usuario
// @Description Request para registro de nuevo usuario
type RegisterRequest struct {
	// Email del usuario (único)
	Email string `json:"email" validate:"required,email,max=255" example:"user@example.com"`
	// Contraseña (mínimo 8 caracteres)
	Password string `json:"password" validate:"required,min=8,max=100" example:"SecurePass123!"`
	// Nombre del usuario
	FirstName string `json:"first_name" validate:"required,min=2,max=100" example:"John"`
	// Apellido del usuario
	LastName string `json:"last_name" validate:"required,min=2,max=100" example:"Doe"`
	// Rol del usuario
	Role Role `json:"role" validate:"required,oneof=SUPER_ADMIN COMPANY CONSUMER" enums:"SUPER_ADMIN,COMPANY,CONSUMER" example:"COMPANY"`
} // @name RegisterRequest

// LoginRequest DTO para inicio de sesión
// @Description Request para autenticación de usuario
type LoginRequest struct {
	// Email del usuario
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
	// Contraseña del usuario
	Password string `json:"password" validate:"required" example:"SecurePass123!"`
} // @name LoginRequest

// RefreshTokenRequest DTO para refrescar token
// @Description Request para renovar access token usando refresh token
type RefreshTokenRequest struct {
	// Refresh token válido
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
} // @name RefreshTokenRequest

// RequestPasswordResetRequest DTO para solicitar reset de contraseña
// @Description Request para solicitar recuperación de contraseña
type RequestPasswordResetRequest struct {
	// Email del usuario
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
} // @name RequestPasswordResetRequest

// ConfirmPasswordResetRequest DTO para confirmar reset de contraseña
// @Description Request para confirmar recuperación de contraseña con token
type ConfirmPasswordResetRequest struct {
	// Token de recuperación recibido por email
	Token string `json:"token" validate:"required" example:"abc123def456"`
	// Nueva contraseña (mínimo 8 caracteres)
	NewPassword string `json:"new_password" validate:"required,min=8,max=100" example:"NewSecurePass123!"`
} // @name ConfirmPasswordResetRequest

// AuthResponse DTO para respuesta de autenticación exitosa
// @Description Respuesta de autenticación exitosa con tokens
type AuthResponse struct {
	// Access token JWT
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// Refresh token JWT
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// Tipo de token (siempre Bearer)
	TokenType string `json:"token_type" example:"Bearer"`
	// Tiempo de expiración en segundos
	ExpiresIn int `json:"expires_in" example:"900"`
	// Datos del usuario autenticado
	User UserDTO `json:"user"`
} // @name AuthResponse

// RefreshTokenResponse DTO para respuesta de refresh token
// @Description Respuesta con nuevo access token
type RefreshTokenResponse struct {
	// Nuevo access token JWT
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// Tipo de token (siempre Bearer)
	TokenType string `json:"token_type" example:"Bearer"`
	// Tiempo de expiración en segundos
	ExpiresIn int `json:"expires_in" example:"900"`
} // @name RefreshTokenResponse

// UserDTO DTO para representar usuario sin información sensible
// @Description Información del usuario (sin datos sensibles)
type UserDTO struct {
	// ID único del usuario
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Email del usuario
	Email string `json:"email" example:"user@example.com"`
	// Nombre del usuario
	FirstName string `json:"first_name" example:"John"`
	// Apellido del usuario
	LastName string `json:"last_name" example:"Doe"`
	// Rol del usuario
	Role Role `json:"role" enums:"SUPER_ADMIN,COMPANY,CONSUMER" example:"COMPANY"`
	// Indica si el usuario está activo
	IsActive bool `json:"is_active" example:"true"`
	// Indica si el email está verificado
	IsEmailVerified bool `json:"is_email_verified" example:"false"`
	// Fecha y hora del último login
	LastLoginAt *time.Time `json:"last_login_at,omitempty" example:"2026-03-02T15:04:05Z"`
	// Fecha de creación
	CreatedAt time.Time `json:"created_at" example:"2026-03-01T10:00:00Z"`
	// Fecha de última actualización
	UpdatedAt time.Time `json:"updated_at" example:"2026-03-02T14:30:00Z"`
} // @name UserDTO

// MessageResponse DTO para respuestas de mensajes simples
// @Description Respuesta con mensaje simple de éxito/error
type MessageResponse struct {
	// Indica si la operación fue exitosa
	Success bool `json:"success" example:"true"`
	// Mensaje descriptivo
	Message string `json:"message" example:"Operation completed successfully"`
} // @name MessageResponse

// ToUserDTO convierte User entity a UserDTO
func ToUserDTO(user *User) UserDTO {
	return UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

// ValidationError representa un error de validación
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse DTO para errores de validación
type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Errors  []ValidationError `json:"errors"`
}
