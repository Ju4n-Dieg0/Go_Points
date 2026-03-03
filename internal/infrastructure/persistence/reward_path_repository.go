package persistence

import (
	"context"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/reward"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// rewardPathRepository implementa reward.RewardPathRepository
type rewardPathRepository struct {
	db *gorm.DB
}

// NewRewardPathRepository crea una nueva instancia
func NewRewardPathRepository(db *gorm.DB) reward.RewardPathRepository {
	return &rewardPathRepository{db: db}
}

// Create crea un nuevo camino
func (r *rewardPathRepository) Create(ctx context.Context, path *reward.RewardPath) error {
	return r.db.WithContext(ctx).Create(path).Error
}

// FindByID busca un camino por ID
func (r *rewardPathRepository) FindByID(ctx context.Context, id uuid.UUID) (*reward.RewardPath, error) {
	var path reward.RewardPath
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&path).Error
	if err != nil {
		return nil, err
	}
	return &path, nil
}

// FindByIDAndCompany busca un camino por ID y CompanyID
func (r *rewardPathRepository) FindByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*reward.RewardPath, error) {
	var path reward.RewardPath
	err := r.db.WithContext(ctx).
		Where("id = ? AND company_id = ?", id, companyID).
		First(&path).Error
	if err != nil {
		return nil, err
	}
	return &path, nil
}

// Update actualiza un camino
func (r *rewardPathRepository) Update(ctx context.Context, path *reward.RewardPath) error {
	return r.db.WithContext(ctx).Save(path).Error
}

// Delete elimina un camino (soft delete)
func (r *rewardPathRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&reward.RewardPath{}, "id = ?", id).Error
}

// ListByCompany lista caminos de una empresa con paginación
func (r *rewardPathRepository) ListByCompany(ctx context.Context, companyID uuid.UUID, page, pageSize int) ([]reward.RewardPath, int64, error) {
	var paths []reward.RewardPath
	var total int64

	// Contar total
	if err := r.db.WithContext(ctx).
		Model(&reward.RewardPath{}).
		Where("company_id = ?", companyID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener caminos paginados
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&paths).Error

	if err != nil {
		return nil, 0, err
	}

	return paths, total, nil
}
