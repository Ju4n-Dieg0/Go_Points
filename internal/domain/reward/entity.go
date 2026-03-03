package reward

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Reward representa una recompensa canjeable por puntos
type Reward struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	ProductID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	RequiredPoints int64          `gorm:"not null" json:"required_points"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Reward) TableName() string {
	return "rewards"
}

// RewardPath representa un camino/colección de recompensas
type RewardPath struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (RewardPath) TableName() string {
	return "reward_paths"
}

// RewardPathItem representa un ítem en un camino de recompensas
type RewardPathItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RewardPathID  uuid.UUID      `gorm:"type:uuid;not null;index:idx_path_order" json:"reward_path_id"`
	RewardID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"reward_id"`
	Order         int            `gorm:"not null;index:idx_path_order" json:"order"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (RewardPathItem) TableName() string {
	return "reward_path_items"
}
