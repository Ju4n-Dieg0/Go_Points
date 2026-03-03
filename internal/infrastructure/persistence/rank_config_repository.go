package persistence

import (
	"context"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// rankConfigRepository implementa point.RankConfigRepository
type rankConfigRepository struct {
	db *gorm.DB
}

// NewRankConfigRepository crea una nueva instancia
func NewRankConfigRepository(db *gorm.DB) point.RankConfigRepository {
	return &rankConfigRepository{db: db}
}

// FindByCompany busca la configuración de una empresa
func (r *rankConfigRepository) FindByCompany(ctx context.Context, companyID uuid.UUID) (*point.CompanyRankConfig, error) {
	var config point.CompanyRankConfig
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Create crea una nueva configuración
func (r *rankConfigRepository) Create(ctx context.Context, config *point.CompanyRankConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// Update actualiza una configuración
func (r *rankConfigRepository) Update(ctx context.Context, config *point.CompanyRankConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// Upsert crea o actualiza una configuración
func (r *rankConfigRepository) Upsert(ctx context.Context, config *point.CompanyRankConfig) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "company_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"silver_min_points", "gold_min_points", "updated_at"}),
		}).
		Create(config).Error
}
