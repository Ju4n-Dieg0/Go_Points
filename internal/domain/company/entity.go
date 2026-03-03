package company

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Company representa la entidad de empresa en el dominio
type Company struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `gorm:"type:varchar(255);not null;index"`
	Logo      *string        `gorm:"type:text"`
	IsActive  bool           `gorm:"default:true;not null;index"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName especifica el nombre de la tabla
func (Company) TableName() string {
	return "companies"
}

// Activate activa la empresa
func (c *Company) Activate() {
	c.IsActive = true
}

// Deactivate desactiva la empresa
func (c *Company) Deactivate() {
	c.IsActive = false
}

// CanOperate verifica si la empresa puede realizar operaciones
func (c *Company) CanOperate() bool {
	return c.IsActive && c.DeletedAt.Time.IsZero()
}
