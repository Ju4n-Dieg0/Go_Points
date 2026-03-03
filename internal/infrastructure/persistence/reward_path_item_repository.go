package persistence

import (
	"context"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/reward"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// rewardPathItemRepository implementa reward.RewardPathItemRepository
type rewardPathItemRepository struct {
	db *gorm.DB
}

// NewRewardPathItemRepository crea una nueva instancia
func NewRewardPathItemRepository(db *gorm.DB) reward.RewardPathItemRepository {
	return &rewardPathItemRepository{db: db}
}

// Create crea un nuevo item
func (r *rewardPathItemRepository) Create(ctx context.Context, item *reward.RewardPathItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// CreateBatch crea múltiples items en lote
func (r *rewardPathItemRepository) CreateBatch(ctx context.Context, items []reward.RewardPathItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

// FindByPath obtiene todos los items de un camino ordenados
func (r *rewardPathItemRepository) FindByPath(ctx context.Context, pathID uuid.UUID) ([]reward.RewardPathItem, error) {
	var items []reward.RewardPathItem
	err := r.db.WithContext(ctx).
		Where("reward_path_id = ?", pathID).
		Order("`order` ASC"). // Backticks porque "order" es palabra reservada
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteByPath elimina todos los items de un camino
func (r *rewardPathItemRepository) DeleteByPath(ctx context.Context, pathID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("reward_path_id = ?", pathID).
		Delete(&reward.RewardPathItem{}).Error
}

// Update actualiza un item
func (r *rewardPathItemRepository) Update(ctx context.Context, item *reward.RewardPathItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}
