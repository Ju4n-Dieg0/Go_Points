package reward

import (
	"context"

	"github.com/google/uuid"
)

// RewardRepository define operaciones de recompensas
type RewardRepository interface {
	// Create crea una nueva recompensa
	Create(ctx context.Context, reward *Reward) error

	// FindByID busca una recompensa por ID
	FindByID(ctx context.Context, id uuid.UUID) (*Reward, error)

	// FindByIDAndCompany busca una recompensa por ID y CompanyID
	FindByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*Reward, error)

	// Update actualiza una recompensa
	Update(ctx context.Context, reward *Reward) error

	// Delete elimina una recompensa (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByCompany lista recompensas de una empresa con paginación
	ListByCompany(ctx context.Context, companyID uuid.UUID, page, pageSize int) ([]Reward, int64, error)

	// FindByProduct busca recompensas por producto
	FindByProduct(ctx context.Context, productID uuid.UUID) ([]Reward, error)
}

// RewardPathRepository define operaciones de caminos de recompensas
type RewardPathRepository interface {
	// Create crea un nuevo camino
	Create(ctx context.Context, path *RewardPath) error

	// FindByID busca un camino por ID
	FindByID(ctx context.Context, id uuid.UUID) (*RewardPath, error)

	// FindByIDAndCompany busca un camino por ID y CompanyID
	FindByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*RewardPath, error)

	// Update actualiza un camino
	Update(ctx context.Context, path *RewardPath) error

	// Delete elimina un camino (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByCompany lista caminos de una empresa con paginación
	ListByCompany(ctx context.Context, companyID uuid.UUID, page, pageSize int) ([]RewardPath, int64, error)
}

// RewardPathItemRepository define operaciones de items de caminos
type RewardPathItemRepository interface {
	// Create crea un nuevo item
	Create(ctx context.Context, item *RewardPathItem) error

	// CreateBatch crea múltiples items
	CreateBatch(ctx context.Context, items []RewardPathItem) error

	// FindByPath lista items de un camino ordenados
	FindByPath(ctx context.Context, pathID uuid.UUID) ([]RewardPathItem, error)

	// DeleteByPath elimina todos los items de un camino
	DeleteByPath(ctx context.Context, pathID uuid.UUID) error

	// Update actualiza un item
	Update(ctx context.Context, item *RewardPathItem) error
}
