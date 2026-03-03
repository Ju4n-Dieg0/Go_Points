package auth

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest DTO para registro de usuario
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Role      Role   `json:"role" validate:"required,oneof=SUPER_ADMIN COMPANY CONSUMER"`
}

// LoginRequest DTO para inicio de sesión
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest DTO para refrescar token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RequestPasswordResetRequest DTO para solicitar reset de contraseña
type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ConfirmPasswordResetRequest DTO para confirmar reset de contraseña
type ConfirmPasswordResetRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

// AuthResponse DTO para respuesta de autenticación exitosa
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	User         UserDTO   `json:"user"`
}

// RefreshTokenResponse DTO para respuesta de refresh token
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// UserDTO DTO para representar usuario sin información sensible
type UserDTO struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Role            Role       `json:"role"`
	IsActive        bool       `json:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// MessageResponse DTO para respuestas de mensajes simples
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

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
