package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product representa un producto de una empresa
type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Photo       string         `gorm:"type:varchar(500)" json:"photo"`
	IsVisible   bool           `gorm:"default:true" json:"is_visible"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (Product) TableName() string {
	return "products"
}
