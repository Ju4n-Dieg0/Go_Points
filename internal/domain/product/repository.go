package product

import (
	"context"

	"github.com/google/uuid"
)

// Repository define la interfaz para operaciones de persistencia de productos
type Repository interface {
	// Create crea un nuevo producto
	Create(ctx context.Context, product *Product) error

	// FindByID busca un producto por ID
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)

	// FindByIDAndCompany busca un producto por ID y CompanyID
	FindByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*Product, error)

	// Update actualiza un producto existente
	Update(ctx context.Context, product *Product) error

	// Delete elimina un producto (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List obtiene productos con paginación filtrados por companyID
	List(ctx context.Context, companyID uuid.UUID, page, pageSize int) ([]Product, int64, error)

	// ListAll obtiene todos los productos visibles con paginación (para consumidores)
	ListAll(ctx context.Context, page, pageSize int) ([]Product, int64, error)

	// Search busca productos por nombre o descripción filtrados por companyID
	Search(ctx context.Context, companyID uuid.UUID, query string, page, pageSize int) ([]Product, int64, error)
}
