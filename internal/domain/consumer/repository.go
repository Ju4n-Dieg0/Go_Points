package consumer

import (
	"context"

	"github.com/google/uuid"
)

// Repository define las operaciones de persistencia para consumidores
type Repository interface {
	// Create crea un nuevo consumidor
	Create(ctx context.Context, consumer *Consumer) error

	// FindByID busca un consumidor por ID
	FindByID(ctx context.Context, id uuid.UUID) (*Consumer, error)

	// FindByDocumentNumber busca un consumidor por número de documento
	FindByDocumentNumber(ctx context.Context, documentNumber string) (*Consumer, error)

	// Update actualiza un consumidor existente
	Update(ctx context.Context, consumer *Consumer) error

	// Delete elimina un consumidor (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List obtiene una lista de consumidores con paginación
	List(ctx context.Context, limit, offset int) ([]*Consumer, error)

	// Count retorna el total de consumidores
	Count(ctx context.Context) (int64, error)

	// ExistsByDocumentNumber verifica si existe un consumidor con el documento dado
	ExistsByDocumentNumber(ctx context.Context, documentNumber string) (bool, error)

	// Search busca consumidores por nombre o email
	Search(ctx context.Context, query string, limit, offset int) ([]*Consumer, error)

	// CountSearch cuenta los resultados de búsqueda
	CountSearch(ctx context.Context, query string) (int64, error)
}
