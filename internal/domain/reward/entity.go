package reward

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Reward representa una recompensa canjeable por puntos
type Reward struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_reward_company;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company_id"`
	ProductID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_reward_product_unique;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"product_id"`
	RequiredPoints int64          `gorm:"not null;check:required_points > 0;index:idx_reward_points" json:"required_points"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index:idx_reward_deleted" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Reward) TableName() string {
	return "rewards"
}

// RewardPath representa un camino/colección de recompensas
type RewardPath struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID uuid.UUID      `gorm:"type:uuid;not null;index:idx_reward_path_company;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company_id"`
	Name      string         `gorm:"type:varchar(255);not null;index:idx_reward_path_name" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_reward_path_deleted" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (RewardPath) TableName() string {
	return "reward_paths"
}

// RewardPathItem representa un ítem en un camino de recompensas
type RewardPathItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RewardPathID  uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_path_reward_unique,priority:1;index:idx_path_order,priority:1;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reward_path_id"`
	RewardID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_path_reward_unique,priority:2;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"reward_id"`
	Order         int            `gorm:"not null;check:\"order\" >= 0;index:idx_path_order,priority:2" json:"order"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_path_item_deleted" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (RewardPathItem) TableName() string {
	return "reward_path_items"
}
