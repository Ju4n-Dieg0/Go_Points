package persistence

import (
	"context"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// balanceRepository implementa point.BalanceRepository
type balanceRepository struct {
	db *gorm.DB
}

// NewBalanceRepository crea una nueva instancia
func NewBalanceRepository(db *gorm.DB) point.BalanceRepository {
	return &balanceRepository{db: db}
}

// FindByConsumerAndCompany busca el balance
func (r *balanceRepository) FindByConsumerAndCompany(ctx context.Context, consumerID, companyID uuid.UUID) (*point.ConsumerCompanyPoints, error) {
	var balance point.ConsumerCompanyPoints
	err := r.db.WithContext(ctx).
		Where("consumer_id = ? AND company_id = ?", consumerID, companyID).
		First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// FindByConsumerAndCompanyForUpdate busca con lock (SELECT FOR UPDATE)
func (r *balanceRepository) FindByConsumerAndCompanyForUpdate(ctx context.Context, tx *gorm.DB, consumerID, companyID uuid.UUID) (*point.ConsumerCompanyPoints, error) {
	var balance point.ConsumerCompanyPoints
	err := tx.WithContext(ctx).
		Raw("SELECT * FROM consumer_company_points WHERE consumer_id = ? AND company_id = ? AND deleted_at IS NULL FOR UPDATE", consumerID, companyID).
		Scan(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// Create crea un nuevo balance
func (r *balanceRepository) Create(ctx context.Context, tx *gorm.DB, balance *point.ConsumerCompanyPoints) error {
	return tx.WithContext(ctx).Create(balance).Error
}

// Update actualiza un balance
func (r *balanceRepository) Update(ctx context.Context, tx *gorm.DB, balance *point.ConsumerCompanyPoints) error {
	return tx.WithContext(ctx).Save(balance).Error
}

// FindInactiveConsumers encuentra consumidores inactivos
func (r *balanceRepository) FindInactiveConsumers(ctx context.Context, inactivityMonths int) ([]point.ConsumerCompanyPoints, error) {
	var balances []point.ConsumerCompanyPoints
	inactiveDate := time.Now().AddDate(0, -inactivityMonths, 0)

	err := r.db.WithContext(ctx).
		Where("last_redemption_date < ? OR last_redemption_date IS NULL", inactiveDate).
		Where("total_available_points > 0").
		Find(&balances).Error

	if err != nil {
		return nil, err
	}
	return balances, nil
}
