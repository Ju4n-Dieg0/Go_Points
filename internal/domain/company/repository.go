package company

import (
	"context"

	"github.com/google/uuid"
)

// Repository define las operaciones de persistencia para empresas
type Repository interface {
	// Create crea una nueva empresa
	Create(ctx context.Context, company *Company) error

	// FindByID busca una empresa por ID
	FindByID(ctx context.Context, id uuid.UUID) (*Company, error)

	// FindByName busca una empresa por nombre
	FindByName(ctx context.Context, name string) (*Company, error)

	// Update actualiza una empresa existente
	Update(ctx context.Context, company *Company) error

	// Delete elimina una empresa (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List obtiene una lista de empresas con paginación
	List(ctx context.Context, limit, offset int, activeOnly bool) ([]*Company, error)

	// Count retorna el total de empresas
	Count(ctx context.Context, activeOnly bool) (int64, error)

	// ExistsByName verifica si existe una empresa con el nombre dado
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Activate activa una empresa
	Activate(ctx context.Context, id uuid.UUID) error

	// Deactivate desactiva una empresa
	Deactivate(ctx context.Context, id uuid.UUID) error

	// FindActiveByID busca una empresa activa por ID
	FindActiveByID(ctx context.Context, id uuid.UUID) (*Company, error)
}
