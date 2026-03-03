package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product representa un producto de una empresa
type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;not null;index:idx_product_company_visible,priority:1;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company_id"`
	Name        string         `gorm:"type:varchar(255);not null;index:idx_product_name" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null;check:price >= 0;index:idx_product_price" json:"price"`
	Photo       string         `gorm:"type:varchar(500)" json:"photo"`
	IsVisible   bool           `gorm:"default:true;not null;index:idx_product_company_visible,priority:2" json:"is_visible"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;not null;index:idx_product_created" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_product_deleted" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (Product) TableName() string {
	return "products"
}
