package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// authRepository implementación concreta del repositorio de autenticación
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository crea una nueva instancia del repositorio de autenticación
func NewAuthRepository(db *gorm.DB) auth.Repository {
	return &authRepository{db: db}
}

// Create crea un nuevo usuario
func (r *authRepository) Create(ctx context.Context, user *auth.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByID busca un usuario por ID
func (r *authRepository) FindByID(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail busca un usuario por email
func (r *authRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByRefreshToken busca un usuario por refresh token
func (r *authRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).
		Where("refresh_token = ? AND refresh_token_expiry > ?", refreshToken, time.Now()).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByPasswordResetToken busca un usuario por token de reset de contraseña
func (r *authRepository) FindByPasswordResetToken(ctx context.Context, token string) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).
		Where("password_reset_token = ? AND password_reset_expiry > ?", token, time.Now()).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update actualiza un usuario existente
func (r *authRepository) Update(ctx context.Context, user *auth.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRefreshToken actualiza el refresh token de un usuario
func (r *authRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken *string, expiry *time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"refresh_token":        refreshToken,
			"refresh_token_expiry": expiry,
		}).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePassword actualiza la contraseña de un usuario
func (r *authRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	if err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":               hashedPassword,
			"password_reset_token":   nil,
			"password_reset_expiry":  nil,
		}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateLastLogin actualiza la última fecha de login
func (r *authRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error; err != nil {
		return err
	}
	return nil
}

// Delete elimina un usuario (soft delete)
func (r *authRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&auth.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *authRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List obtiene una lista de usuarios con paginación
func (r *authRepository) List(ctx context.Context, limit, offset int) ([]*auth.User, error) {
	var users []*auth.User
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Count retorna el total de usuarios
func (r *authRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&auth.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
