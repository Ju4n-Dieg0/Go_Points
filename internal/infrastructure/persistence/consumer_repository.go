package persistence

import (
	"context"
	"errors"
	"strings"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/consumer"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// consumerRepository implementación concreta del repositorio de consumidores
type consumerRepository struct {
	db *gorm.DB
}

// NewConsumerRepository crea una nueva instancia del repositorio de consumidores
func NewConsumerRepository(db *gorm.DB) consumer.Repository {
	return &consumerRepository{db: db}
}

// Create crea un nuevo consumidor
func (r *consumerRepository) Create(ctx context.Context, cons *consumer.Consumer) error {
	if err := r.db.WithContext(ctx).Create(cons).Error; err != nil {
		// Manejar error de unique violation específicamente
		if isUniqueViolation(err) {
			return &UniqueViolationError{
				Field:   "document_number",
				Value:   cons.DocumentNumber,
				Message: "consumer with this document number already exists",
			}
		}
		return err
	}
	return nil
}

// FindByID busca un consumidor por ID
func (r *consumerRepository) FindByID(ctx context.Context, id uuid.UUID) (*consumer.Consumer, error) {
	var cons consumer.Consumer
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&cons).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cons, nil
}

// FindByDocumentNumber busca un consumidor por número de documento
func (r *consumerRepository) FindByDocumentNumber(ctx context.Context, documentNumber string) (*consumer.Consumer, error) {
	var cons consumer.Consumer
	if err := r.db.WithContext(ctx).Where("document_number = ?", documentNumber).First(&cons).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cons, nil
}

// Update actualiza un consumidor existente
func (r *consumerRepository) Update(ctx context.Context, cons *consumer.Consumer) error {
	if err := r.db.WithContext(ctx).Save(cons).Error; err != nil {
		if isUniqueViolation(err) {
			return &UniqueViolationError{
				Field:   "document_number",
				Value:   cons.DocumentNumber,
				Message: "consumer with this document number already exists",
			}
		}
		return err
	}
	return nil
}

// Delete elimina un consumidor (soft delete)
func (r *consumerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&consumer.Consumer{}, id).Error; err != nil {
		return err
	}
	return nil
}

// List obtiene una lista de consumidores con paginación
func (r *consumerRepository) List(ctx context.Context, limit, offset int) ([]*consumer.Consumer, error) {
	var consumers []*consumer.Consumer
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&consumers).Error; err != nil {
		return nil, err
	}
	return consumers, nil
}

// Count retorna el total de consumidores
func (r *consumerRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&consumer.Consumer{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByDocumentNumber verifica si existe un consumidor con el documento dado
func (r *consumerRepository) ExistsByDocumentNumber(ctx context.Context, documentNumber string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&consumer.Consumer{}).
		Where("document_number = ?", documentNumber).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Search busca consumidores por nombre o email
func (r *consumerRepository) Search(ctx context.Context, query string, limit, offset int) ([]*consumer.Consumer, error) {
	var consumers []*consumer.Consumer
	searchPattern := "%" + strings.ToLower(query) + "%"

	if err := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR document_number LIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&consumers).Error; err != nil {
		return nil, err
	}
	return consumers, nil
}

// CountSearch cuenta los resultados de búsqueda
func (r *consumerRepository) CountSearch(ctx context.Context, query string) (int64, error) {
	var count int64
	searchPattern := "%" + strings.ToLower(query) + "%"

	if err := r.db.WithContext(ctx).
		Model(&consumer.Consumer{}).
		Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR document_number LIKE ?", searchPattern, searchPattern, searchPattern).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UniqueViolationError error personalizado para violaciones de constraint único
type UniqueViolationError struct {
	Field   string
	Value   string
	Message string
}

func (e *UniqueViolationError) Error() string {
	return e.Message
}

// isUniqueViolation verifica si un error es una violación de constraint único
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	// PostgreSQL error code 23505 es unique violation
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
