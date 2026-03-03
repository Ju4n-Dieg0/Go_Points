package point

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionType define los tipos de transacciones de puntos
type TransactionType string

const (
	TransactionTypeEarn    TransactionType = "EARN"    // Ganar puntos
	TransactionTypeRedeem  TransactionType = "REDEEM"  // Redimir puntos
	TransactionTypeExpire  TransactionType = "EXPIRE"  // Expiración automática
	TransactionTypePenalty TransactionType = "PENALTY" // Penalización por inactividad
)

// Rank define los rangos de clientes
type Rank string

const (
	RankBronze Rank = "BRONZE" // 0 puntos
	RankSilver Rank = "SILVER" // Configurable
	RankGold   Rank = "GOLD"   // Configurable
)

// ConsumerCompanyPoints representa el balance de puntos de un consumidor en una empresa
type ConsumerCompanyPoints struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConsumerID             uuid.UUID      `gorm:"type:uuid;not null;index:idx_consumer_company,unique" json:"consumer_id"`
	CompanyID              uuid.UUID      `gorm:"type:uuid;not null;index:idx_consumer_company,unique" json:"company_id"`
	TotalHistoricalPoints  int64          `gorm:"not null;default:0" json:"total_historical_points"`
	TotalAvailablePoints   int64          `gorm:"not null;default:0" json:"total_available_points"`
	LastRedemptionDate     *time.Time     `gorm:"index" json:"last_redemption_date,omitempty"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (ConsumerCompanyPoints) TableName() string {
	return "consumer_company_points"
}

// GetRank calcula el rango actual basado en puntos históricos
func (c *ConsumerCompanyPoints) GetRank(config *CompanyRankConfig) Rank {
	if config == nil {
		return RankBronze
	}

	if c.TotalHistoricalPoints >= config.GoldMinPoints {
		return RankGold
	}
	if c.TotalHistoricalPoints >= config.SilverMinPoints {
		return RankSilver
	}
	return RankBronze
}

// PointTransaction representa una transacción de puntos
type PointTransaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConsumerID      uuid.UUID       `gorm:"type:uuid;not null;index:idx_consumer_transactions" json:"consumer_id"`
	CompanyID       uuid.UUID       `gorm:"type:uuid;not null;index:idx_company_transactions" json:"company_id"`
	Points          int64           `gorm:"not null" json:"points"`
	RemainingPoints int64           `gorm:"not null;default:0" json:"remaining_points"`
	Type            TransactionType `gorm:"type:varchar(20);not null;index" json:"type"`
	ExpirationDate  *time.Time      `gorm:"index" json:"expiration_date,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName especifica el nombre de la tabla
func (PointTransaction) TableName() string {
	return "point_transactions"
}

// IsExpired verifica si la transacción ha expirado
func (t *PointTransaction) IsExpired() bool {
	if t.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*t.ExpirationDate)
}

// HasAvailablePoints verifica si la transacción tiene puntos disponibles
func (t *PointTransaction) HasAvailablePoints() bool {
	return t.RemainingPoints > 0
}

// CompanyRankConfig define la configuración de rangos para una empresa
type CompanyRankConfig struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"company_id"`
	SilverMinPoints int64         `gorm:"not null;default:1000" json:"silver_min_points"`
	GoldMinPoints   int64         `gorm:"not null;default:5000" json:"gold_min_points"`
	CreatedAt       time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (CompanyRankConfig) TableName() string {
	return "company_rank_configs"
}

// Validate valida que la configuración sea válida
func (c *CompanyRankConfig) Validate() error {
	if c.SilverMinPoints <= 0 {
		return ErrInvalidRankConfig
	}
	if c.GoldMinPoints <= c.SilverMinPoints {
		return ErrInvalidRankConfig
	}
	return nil
}

// ErrInvalidRankConfig error cuando la configuración de rangos es inválida
var ErrInvalidRankConfig = &ValidationError{Message: "gold_min_points must be greater than silver_min_points"}

// ValidationError representa un error de validación
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
