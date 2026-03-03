package persistence

import (
	"context"
	"errors"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/company"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// companyRepository implementación concreta del repositorio de empresas
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository crea una nueva instancia del repositorio de empresas
func NewCompanyRepository(db *gorm.DB) company.Repository {
	return &companyRepository{db: db}
}

// Create crea una nueva empresa
func (r *companyRepository) Create(ctx context.Context, comp *company.Company) error {
	if err := r.db.WithContext(ctx).Create(comp).Error; err != nil {
		return err
	}
	return nil
}

// FindByID busca una empresa por ID
func (r *companyRepository) FindByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	var comp company.Company
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&comp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comp, nil
}

// FindByName busca una empresa por nombre
func (r *companyRepository) FindByName(ctx context.Context, name string) (*company.Company, error) {
	var comp company.Company
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&comp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comp, nil
}

// Update actualiza una empresa existente
func (r *companyRepository) Update(ctx context.Context, comp *company.Company) error {
	if err := r.db.WithContext(ctx).Save(comp).Error; err != nil {
		return err
	}
	return nil
}

// Delete elimina una empresa (soft delete)
func (r *companyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&company.Company{}, id).Error; err != nil {
		return err
	}
	return nil
}

// List obtiene una lista de empresas con paginación
func (r *companyRepository) List(ctx context.Context, limit, offset int, activeOnly bool) ([]*company.Company, error) {
	var companies []*company.Company
	query := r.db.WithContext(ctx)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

// Count retorna el total de empresas
func (r *companyRepository) Count(ctx context.Context, activeOnly bool) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&company.Company{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByName verifica si existe una empresa con el nombre dado
func (r *companyRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&company.Company{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Activate activa una empresa
func (r *companyRepository) Activate(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&company.Company{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

// Deactivate desactiva una empresa
func (r *companyRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&company.Company{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

// FindActiveByID busca una empresa activa por ID
func (r *companyRepository) FindActiveByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	var comp company.Company
	if err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&comp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comp, nil
}
