package persistence

import (
	"context"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// productRepository implementa product.Repository usando GORM
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository crea una nueva instancia de productRepository
func NewProductRepository(db *gorm.DB) product.Repository {
	return &productRepository{db: db}
}

// Create crea un nuevo producto
func (r *productRepository) Create(ctx context.Context, prod *product.Product) error {
	return r.db.WithContext(ctx).Create(prod).Error
}

// FindByID busca un producto por ID
func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	var prod product.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&prod).Error
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

// FindByIDAndCompany busca un producto por ID y CompanyID
func (r *productRepository) FindByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*product.Product, error) {
	var prod product.Product
	err := r.db.WithContext(ctx).
		Where("id = ? AND company_id = ?", id, companyID).
		First(&prod).Error
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

// Update actualiza un producto existente
func (r *productRepository) Update(ctx context.Context, prod *product.Product) error {
	return r.db.WithContext(ctx).Save(prod).Error
}

// Delete elimina un producto (soft delete)
func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&product.Product{}, "id = ?", id).Error
}

// List obtiene productos con paginación filtrados por companyID
func (r *productRepository) List(ctx context.Context, companyID uuid.UUID, page, pageSize int) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	// Contar total
	if err := r.db.WithContext(ctx).
		Model(&product.Product{}).
		Where("company_id = ?", companyID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener productos paginados
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListAll obtiene todos los productos visibles con paginación
func (r *productRepository) ListAll(ctx context.Context, page, pageSize int) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	// Contar total de productos visibles
	if err := r.db.WithContext(ctx).
		Model(&product.Product{}).
		Where("is_visible = ?", true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener productos visibles paginados
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("is_visible = ?", true).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Search busca productos por nombre o descripción filtrados por companyID
func (r *productRepository) Search(ctx context.Context, companyID uuid.UUID, query string, page, pageSize int) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	searchPattern := "%" + query + "%"

	// Contar total
	if err := r.db.WithContext(ctx).
		Model(&product.Product{}).
		Where("company_id = ? AND (name ILIKE ? OR description ILIKE ?)", companyID, searchPattern, searchPattern).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener productos paginados
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND (name ILIKE ? OR description ILIKE ?)", companyID, searchPattern, searchPattern).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
