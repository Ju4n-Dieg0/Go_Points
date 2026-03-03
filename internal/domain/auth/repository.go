package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository define las operaciones de persistencia para usuarios
type Repository interface {
	// Create crea un nuevo usuario
	Create(ctx context.Context, user *User) error

	// FindByID busca un usuario por ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail busca un usuario por email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByRefreshToken busca un usuario por refresh token
	FindByRefreshToken(ctx context.Context, refreshToken string) (*User, error)

	// FindByPasswordResetToken busca un usuario por token de reset de contraseña
	FindByPasswordResetToken(ctx context.Context, token string) (*User, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *User) error

	// UpdateRefreshToken actualiza el refresh token de un usuario
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken *string, expiry *time.Time) error

	// UpdatePassword actualiza la contraseña de un usuario
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error

	// UpdateLastLogin actualiza la última fecha de login
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Delete elimina un usuario (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByEmail verifica si existe un usuario con el email dado
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// List obtiene una lista de usuarios con paginación
	List(ctx context.Context, limit, offset int) ([]*User, error)

	// Count retorna el total de usuarios
	Count(ctx context.Context) (int64, error)
}
