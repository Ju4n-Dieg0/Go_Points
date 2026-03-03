package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role representa los roles de usuario en el sistema
type Role string

const (
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleCompany    Role = "COMPANY"
	RoleConsumer   Role = "CONSUMER"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email                string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password             string         `gorm:"type:varchar(255);not null"`
	FirstName            string         `gorm:"type:varchar(100);not null"`
	LastName             string         `gorm:"type:varchar(100);not null"`
	Role                 Role           `gorm:"type:varchar(50);not null;default:'CONSUMER'"`
	IsActive             bool           `gorm:"default:true;not null"`
	IsEmailVerified      bool           `gorm:"default:false;not null"`
	EmailVerificationToken *string      `gorm:"type:varchar(255)"`
	PasswordResetToken   *string        `gorm:"type:varchar(255)"`
	PasswordResetExpiry  *time.Time
	RefreshToken         *string        `gorm:"type:text"`
	RefreshTokenExpiry   *time.Time
	LastLoginAt          *time.Time
	CreatedAt            time.Time      `gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "users"
}

// IsValidRole verifica si el rol es válido
func IsValidRole(role Role) bool {
	switch role {
	case RoleSuperAdmin, RoleCompany, RoleConsumer:
		return true
	default:
		return false
	}
}

// HasRole verifica si el usuario tiene un rol específico
func (u *User) HasRole(role Role) bool {
	return u.Role == role
}

// IsSuperAdmin verifica si el usuario es super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsCompany verifica si el usuario es empresa
func (u *User) IsCompany() bool {
	return u.Role == RoleCompany
}

// IsConsumer verifica si el usuario es consumidor
func (u *User) IsConsumer() bool {
	return u.Role == RoleConsumer
}

// CanAccessResource verifica si el usuario puede acceder a un recurso
func (u *User) CanAccessResource(requiredRoles []Role) bool {
	for _, role := range requiredRoles {
		if u.Role == role {
			return true
		}
	}
	return false
}
