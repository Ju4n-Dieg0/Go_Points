package subscription

import (
	"context"
	"fmt"

	domainCompany "github.com/Ju4n-Dieg0/Go_Points/internal/domain/company"
	domainSubscription "github.com/Ju4n-Dieg0/Go_Points/internal/domain/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service define la lógica de negocio de suscripciones
type Service interface {
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*domainSubscription.SubscriptionResponse, error)
	Renew(ctx context.Context, companyID uuid.UUID) (*domainSubscription.SubscriptionResponse, error)
	Cancel(ctx context.Context, companyID uuid.UUID) error
	CheckAndDeactivateExpired(ctx context.Context) error
	ValidateActiveSubscription(ctx context.Context, companyID uuid.UUID) error
}

type service struct {
	subscriptionRepo domainSubscription.Repository
	companyRepo      domainCompany.Repository
	db               *gorm.DB
}

// NewService crea una nueva instancia del servicio de suscripciones
func NewService(
	subscriptionRepo domainSubscription.Repository,
	companyRepo domainCompany.Repository,
	db *gorm.DB,
) Service {
	return &service{
		subscriptionRepo: subscriptionRepo,
		companyRepo:      companyRepo,
		db:               db,
	}
}

// GetByCompanyID obtiene la suscripción de una empresa
func (s *service) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*domainSubscription.SubscriptionResponse, error) {
	// Verificar que la empresa existe
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	// Obtener suscripción
	subscription, err := s.subscriptionRepo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if subscription == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("subscription not found"))
	}

	response := domainSubscription.ToSubscriptionResponse(subscription)
	return &response, nil
}

// Renew renueva la suscripción de una empresa en una transacción
func (s *service) Renew(ctx context.Context, companyID uuid.UUID) (*domainSubscription.SubscriptionResponse, error) {
	// Verificar que la empresa existe
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	var subscription *domainSubscription.Subscription

	// Renovar en una transacción
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Obtener suscripción actual
		sub, err := s.subscriptionRepo.FindByCompanyID(ctx, companyID)
		if err != nil {
			return err
		}
		if sub == nil {
			return fmt.Errorf("subscription not found")
		}

		subscription = sub

		// Renovar suscripción (suma 30 días)
		subscription.Renew()

		// Actualizar suscripción
		if err := tx.Save(subscription).Error; err != nil {
			return err
		}

		// Si la empresa estaba inactiva, activarla
		if !company.IsActive {
			if err := tx.Model(&domainCompany.Company{}).
				Where("id = ?", companyID).
				Update("is_active", true).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainSubscription.ToSubscriptionResponse(subscription)
	return &response, nil
}

// Cancel cancela la suscripción de una empresa (no la elimina, solo deja que expire)
func (s *service) Cancel(ctx context.Context, companyID uuid.UUID) error {
	// Verificar que la empresa existe
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	// Obtener suscripción
	subscription, err := s.subscriptionRepo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if subscription == nil {
		return errors.ErrNotFound.WithError(fmt.Errorf("subscription not found"))
	}

	// Cancelar suscripción
	subscription.Cancel()

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// CheckAndDeactivateExpired verifica y desactiva suscripciones expiradas
// Este método debe ser llamado periódicamente (ej: cronjob)
func (s *service) CheckAndDeactivateExpired(ctx context.Context) error {
	// Obtener suscripciones expiradas
	expiredSubs, err := s.subscriptionRepo.FindExpiredSubscriptions(ctx)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	if len(expiredSubs) == 0 {
		return nil
	}

	// Procesar cada suscripción expirada en una transacción
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Desactivar suscripciones expiradas
		if err := tx.Model(&domainSubscription.Subscription{}).
			Where("is_active = ? AND end_date < NOW()", true).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Desactivar empresas con suscripciones expiradas
		companyIDs := make([]uuid.UUID, 0, len(expiredSubs))
		for _, sub := range expiredSubs {
			companyIDs = append(companyIDs, sub.CompanyID)
		}

		if len(companyIDs) > 0 {
			if err := tx.Model(&domainCompany.Company{}).
				Where("id IN ? AND is_active = ?", companyIDs, true).
				Update("is_active", false).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ValidateActiveSubscription valida que una empresa tenga suscripción activa
func (s *service) ValidateActiveSubscription(ctx context.Context, companyID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.FindActiveByCompanyID(ctx, companyID)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	if subscription == nil {
		return errors.ErrForbidden.WithError(fmt.Errorf("no active subscription found"))
	}

	if subscription.IsExpired() {
		// Desactivar empresa y suscripción
		_ = s.companyRepo.Deactivate(ctx, companyID)
		subscription.Cancel()
		_ = s.subscriptionRepo.Update(ctx, subscription)

		return errors.ErrForbidden.WithError(fmt.Errorf("subscription has expired"))
	}

	return nil
}
