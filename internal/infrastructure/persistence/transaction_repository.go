package persistence

import (
	"context"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// transactionRepository implementa point.TransactionRepository
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository crea una nueva instancia
func NewTransactionRepository(db *gorm.DB) point.TransactionRepository {
	return &transactionRepository{db: db}
}

// Create crea una nueva transacción
func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *point.PointTransaction) error {
	return tx.WithContext(ctx).Create(transaction).Error
}

// FindAvailableForRedemption encuentra transacciones disponibles en orden FIFO
// CRÍTICO: Usa SELECT FOR UPDATE para evitar race conditions
func (r *transactionRepository) FindAvailableForRedemption(ctx context.Context, tx *gorm.DB, consumerID, companyID uuid.UUID) ([]point.PointTransaction, error) {
	var transactions []point.PointTransaction

	err := tx.WithContext(ctx).
		Raw(`SELECT * FROM point_transactions 
			WHERE consumer_id = ? 
			AND company_id = ? 
			AND type = ? 
			AND remaining_points > 0 
			AND (expiration_date IS NULL OR expiration_date > ?) 
			ORDER BY created_at ASC 
			FOR UPDATE`,
			consumerID, companyID, point.TransactionTypeEarn, time.Now()).
		Scan(&transactions).Error

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// UpdateBatch actualiza múltiples transacciones
func (r *transactionRepository) UpdateBatch(ctx context.Context, tx *gorm.DB, transactions []point.PointTransaction) error {
	for i := range transactions {
		if err := tx.WithContext(ctx).Save(&transactions[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

// FindByConsumerAndCompany lista transacciones con paginación
func (r *transactionRepository) FindByConsumerAndCompany(ctx context.Context, consumerID, companyID uuid.UUID, page, pageSize int) ([]point.PointTransaction, int64, error) {
	var transactions []point.PointTransaction
	var total int64

	// Contar total
	if err := r.db.WithContext(ctx).
		Model(&point.PointTransaction{}).
		Where("consumer_id = ? AND company_id = ?", consumerID, companyID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener transacciones paginadas
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("consumer_id = ? AND company_id = ?", consumerID, companyID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindExpired encuentra transacciones expiradas con puntos disponibles
func (r *transactionRepository) FindExpired(ctx context.Context) ([]point.PointTransaction, error) {
	var transactions []point.PointTransaction

	err := r.db.WithContext(ctx).
		Where("type = ?", point.TransactionTypeEarn).
		Where("remaining_points > 0").
		Where("expiration_date IS NOT NULL").
		Where("expiration_date < ?", time.Now()).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// FindExpiringSoon encuentra transacciones que expirarán pronto
func (r *transactionRepository) FindExpiringSoon(ctx context.Context, daysBeforeExpiration int) ([]point.PointTransaction, error) {
	var transactions []point.PointTransaction

	thresholdDate := time.Now().AddDate(0, 0, daysBeforeExpiration)

	err := r.db.WithContext(ctx).
		Where("type = ?", point.TransactionTypeEarn).
		Where("remaining_points > 0").
		Where("expiration_date IS NOT NULL").
		Where("expiration_date BETWEEN ? AND ?", time.Now(), thresholdDate).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}
	return transactions, nil
}
