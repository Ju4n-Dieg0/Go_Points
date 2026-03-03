package subscription

import (
	"context"

	"github.com/google/uuid"
)

// Repository define las operaciones de persistencia para suscripciones
type Repository interface {
	// Create crea una nueva suscripción
	Create(ctx context.Context, subscription *Subscription) error

	// FindByID busca una suscripción por ID
	FindByID(ctx context.Context, id uuid.UUID) (*Subscription, error)

	// FindByCompanyID busca la suscripción activa de una empresa
	FindByCompanyID(ctx context.Context, companyID uuid.UUID) (*Subscription, error)

	// FindActiveByCompanyID busca la suscripción activa y no expirada de una empresa
	FindActiveByCompanyID(ctx context.Context, companyID uuid.UUID) (*Subscription, error)

	// Update actualiza una suscripción existente
	Update(ctx context.Context, subscription *Subscription) error

	// List obtiene una lista de suscripciones con paginación
	List(ctx context.Context, limit, offset int) ([]*Subscription, error)

	// Count retorna el total de suscripciones
	Count(ctx context.Context) (int64, error)

	// FindExpiredSubscriptions busca suscripciones expiradas que aún están marcadas como activas
	FindExpiredSubscriptions(ctx context.Context) ([]*Subscription, error)

	// DeactivateExpired desactiva todas las suscripciones expiradas
	DeactivateExpired(ctx context.Context) error

	// ExistsByCompanyID verifica si existe una suscripción para una empresa
	ExistsByCompanyID(ctx context.Context, companyID uuid.UUID) (bool, error)
}
