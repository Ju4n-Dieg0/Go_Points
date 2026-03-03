package persistence

import (
	"context"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/reward"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// rewardRepository implementa reward.RewardRepository
type rewardRepository struct {
	db *gorm.DB
}

// NewRewardRepository crea una nueva instancia
func NewRewardRepository(db *gorm.DB) reward.RewardRepository {
	return &rewardRepository{db: db}
}

// Create crea una nueva recompensa
func (r *rewardRepository) Create(ctx context.Context, rew *reward.Reward) error {
	return r.db.WithContext(ctx).Create(rew).Error
}

// FindByID busca una recompensa por ID
func (r *rewardRepository) FindByID(ctx context.Context, id uuid.UUID) (*reward.Reward, error) {
	var rew reward.Reward
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rew).Error
	if err != nil {
		return nil, err
	}
	return &rew, nil
}

// FindByIDAndCompany busca una recompensa por ID y CompanyID
func (r *rewardRepository) FindByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*reward.Reward, error) {
	var rew reward.Reward
	err := r.db.WithContext(ctx).
		Where("id = ? AND company_id = ?", id, companyID).
		First(&rew).Error
	if err != nil {
		return nil, err
	}
	return &rew, nil
}

// Update actualiza una recompensa
func (r *rewardRepository) Update(ctx context.Context, rew *reward.Reward) error {
	return r.db.WithContext(ctx).Save(rew).Error
}

// Delete elimina una recompensa (soft delete)
func (r *rewardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&reward.Reward{}, "id = ?", id).Error
}

// ListByCompany lista recompensas de una empresa con paginación
func (r *rewardRepository) ListByCompany(ctx context.Context, companyID uuid.UUID, page, pageSize int) ([]reward.Reward, int64, error) {
	var rewards []reward.Reward
	var total int64

	// Contar total
	if err := r.db.WithContext(ctx).
		Model(&reward.Reward{}).
		Where("company_id = ?", companyID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener recompensas paginadas
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rewards).Error

	if err != nil {
		return nil, 0, err
	}

	return rewards, total, nil
}

// FindByProduct busca recompensas por producto
func (r *rewardRepository) FindByProduct(ctx context.Context, productID uuid.UUID) ([]reward.Reward, error) {
	var rewards []reward.Reward
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&rewards).Error
	if err != nil {
		return nil, err
	}
	return rewards, nil
}
