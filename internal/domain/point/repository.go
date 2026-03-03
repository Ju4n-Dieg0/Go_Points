package point

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BalanceRepository define operaciones de balance de puntos
type BalanceRepository interface {
	// FindByConsumerAndCompany busca el balance de un consumidor en una empresa
	FindByConsumerAndCompany(ctx context.Context, consumerID, companyID uuid.UUID) (*ConsumerCompanyPoints, error)

	// FindByConsumerAndCompanyForUpdate busca el balance con lock para actualización
	FindByConsumerAndCompanyForUpdate(ctx context.Context, tx *gorm.DB, consumerID, companyID uuid.UUID) (*ConsumerCompanyPoints, error)

	// Create crea un nuevo balance
	Create(ctx context.Context, tx *gorm.DB, balance *ConsumerCompanyPoints) error

	// Update actualiza un balance
	Update(ctx context.Context, tx *gorm.DB, balance *ConsumerCompanyPoints) error

	// FindInactiveConsumers encuentra consumidores inactivos
	FindInactiveConsumers(ctx context.Context, inactivityMonths int) ([]ConsumerCompanyPoints, error)
}

// TransactionRepository define operaciones de transacciones de puntos
type TransactionRepository interface {
	// Create crea una nueva transacción
	Create(ctx context.Context, tx *gorm.DB, transaction *PointTransaction) error

	// FindAvailableForRedemption encuentra transacciones disponibles para redimir (FIFO)
	FindAvailableForRedemption(ctx context.Context, tx *gorm.DB, consumerID, companyID uuid.UUID) ([]PointTransaction, error)

	// UpdateBatch actualiza múltiples transacciones
	UpdateBatch(ctx context.Context, tx *gorm.DB, transactions []PointTransaction) error

	// FindByConsumerAndCompany lista transacciones con paginación
	FindByConsumerAndCompany(ctx context.Context, consumerID, companyID uuid.UUID, page, pageSize int) ([]PointTransaction, int64, error)

	// FindExpired encuentra transacciones expiradas con puntos disponibles
	FindExpired(ctx context.Context) ([]PointTransaction, error)

	// FindExpiringSoon encuentra transacciones que expirarán pronto
	FindExpiringSoon(ctx context.Context, daysBeforeExpiration int) ([]PointTransaction, error)
}

// RankConfigRepository define operaciones de configuración de rangos
type RankConfigRepository interface {
	// FindByCompany busca la configuración de una empresa
	FindByCompany(ctx context.Context, companyID uuid.UUID) (*CompanyRankConfig, error)

	// Create crea una nueva configuración
	Create(ctx context.Context, config *CompanyRankConfig) error

	// Update actualiza una configuración
	Update(ctx context.Context, config *CompanyRankConfig) error

	// Upsert crea o actualiza una configuración
	Upsert(ctx context.Context, config *CompanyRankConfig) error
}
