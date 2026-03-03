package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/subscription"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// subscriptionRepository implementación concreta del repositorio de suscripciones
type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository crea una nueva instancia del repositorio de suscripciones
func NewSubscriptionRepository(db *gorm.DB) subscription.Repository {
	return &subscriptionRepository{db: db}
}

// Create crea una nueva suscripción
func (r *subscriptionRepository) Create(ctx context.Context, sub *subscription.Subscription) error {
	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return err
	}
	return nil
}

// FindByID busca una suscripción por ID
func (r *subscriptionRepository) FindByID(ctx context.Context, id uuid.UUID) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// FindByCompanyID busca la suscripción activa de una empresa
func (r *subscriptionRepository) FindByCompanyID(ctx context.Context, companyID uuid.UUID) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND is_active = ?", companyID, true).
		Order("created_at DESC").
		First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// FindActiveByCompanyID busca la suscripción activa y no expirada de una empresa
func (r *subscriptionRepository) FindActiveByCompanyID(ctx context.Context, companyID uuid.UUID) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND is_active = ? AND end_date > ?", companyID, true, now).
		Order("created_at DESC").
		First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// Update actualiza una suscripción existente
func (r *subscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {
	if err := r.db.WithContext(ctx).Save(sub).Error; err != nil {
		return err
	}
	return nil
}

// List obtiene una lista de suscripciones con paginación
func (r *subscriptionRepository) List(ctx context.Context, limit, offset int) ([]*subscription.Subscription, error) {
	var subscriptions []*subscription.Subscription
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// Count retorna el total de suscripciones
func (r *subscriptionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&subscription.Subscription{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindExpiredSubscriptions busca suscripciones expiradas que aún están marcadas como activas
func (r *subscriptionRepository) FindExpiredSubscriptions(ctx context.Context) ([]*subscription.Subscription, error) {
	var subscriptions []*subscription.Subscription
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND end_date < ?", true, now).
		Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// DeactivateExpired desactiva todas las suscripciones expiradas
func (r *subscriptionRepository) DeactivateExpired(ctx context.Context) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&subscription.Subscription{}).
		Where("is_active = ? AND end_date < ?", true, now).
		Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByCompanyID verifica si existe una suscripción para una empresa
func (r *subscriptionRepository) ExistsByCompanyID(ctx context.Context, companyID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&subscription.Subscription{}).
		Where("company_id = ?", companyID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
