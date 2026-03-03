package company

import (
	"context"
	"fmt"
	"time"

	domainCompany "github.com/Ju4n-Dieg0/Go_Points/internal/domain/company"
	domainSubscription "github.com/Ju4n-Dieg0/Go_Points/internal/domain/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service define la lógica de negocio de empresas
type Service interface {
	Create(ctx context.Context, req *domainCompany.CreateCompanyRequest) (*domainCompany.CompanyResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainCompany.CompanyResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *domainCompany.UpdateCompanyRequest) (*domainCompany.CompanyResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int, activeOnly bool) (*domainCompany.CompanyListResponse, error)
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type service struct {
	companyRepo      domainCompany.Repository
	subscriptionRepo domainSubscription.Repository
	db               *gorm.DB
}

// NewService crea una nueva instancia del servicio de empresas
func NewService(
	companyRepo domainCompany.Repository,
	subscriptionRepo domainSubscription.Repository,
	db *gorm.DB,
) Service {
	return &service{
		companyRepo:      companyRepo,
		subscriptionRepo: subscriptionRepo,
		db:               db,
	}
}

// Create crea una nueva empresa con su suscripción inicial en una transacción
func (s *service) Create(ctx context.Context, req *domainCompany.CreateCompanyRequest) (*domainCompany.CompanyResponse, error) {
	// Verificar si ya existe una empresa con ese nombre
	exists, err := s.companyRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if exists {
		return nil, errors.ErrConflict.WithError(fmt.Errorf("company with name %s already exists", req.Name))
	}

	// Crear empresa y suscripción en una transacción
	var company *domainCompany.Company
	var subscription *domainSubscription.Subscription

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Crear empresa
		company = &domainCompany.Company{
			ID:       uuid.New(),
			Name:     req.Name,
			Logo:     req.Logo,
			IsActive: true,
		}

		if err := tx.Create(company).Error; err != nil {
			return err
		}

		// Crear suscripción inicial de 30 días
		now := time.Now()
		subscription = &domainSubscription.Subscription{
			ID:        uuid.New(),
			CompanyID: company.ID,
			StartDate: now,
			EndDate:   now.AddDate(0, 0, 30),
			IsActive:  true,
		}

		if err := tx.Create(subscription).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Construir respuesta
	response := domainCompany.ToCompanyResponse(company)
	response.Subscription = &domainCompany.SubscriptionInfo{
		ID:            subscription.ID,
		StartDate:     subscription.StartDate,
		EndDate:       subscription.EndDate,
		IsActive:      subscription.IsActive,
		DaysRemaining: subscription.DaysRemaining(),
	}

	return &response, nil
}

// GetByID obtiene una empresa por ID con su suscripción
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domainCompany.CompanyResponse, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	// Obtener suscripción activa
	subscription, err := s.subscriptionRepo.FindByCompanyID(ctx, company.ID)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainCompany.ToCompanyResponse(company)
	if subscription != nil {
		response.Subscription = &domainCompany.SubscriptionInfo{
			ID:            subscription.ID,
			StartDate:     subscription.StartDate,
			EndDate:       subscription.EndDate,
			IsActive:      subscription.IsActive,
			DaysRemaining: subscription.DaysRemaining(),
		}
	}

	return &response, nil
}

// Update actualiza una empresa existente
func (s *service) Update(ctx context.Context, id uuid.UUID, req *domainCompany.UpdateCompanyRequest) (*domainCompany.CompanyResponse, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	// Actualizar campos
	if req.Name != nil {
		// Verificar si el nuevo nombre ya existe
		if *req.Name != company.Name {
			exists, err := s.companyRepo.ExistsByName(ctx, *req.Name)
			if err != nil {
				return nil, errors.ErrDatabase.WithError(err)
			}
			if exists {
				return nil, errors.ErrConflict.WithError(fmt.Errorf("company with name %s already exists", *req.Name))
			}
		}
		company.Name = *req.Name
	}

	if req.Logo != nil {
		company.Logo = req.Logo
	}

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return s.GetByID(ctx, company.ID)
}

// Delete elimina una empresa (soft delete)
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	if err := s.companyRepo.Delete(ctx, id); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// List obtiene una lista de empresas con paginación
func (s *service) List(ctx context.Context, page, pageSize int, activeOnly bool) (*domainCompany.CompanyListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	companies, err := s.companyRepo.List(ctx, pageSize, offset, activeOnly)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	total, err := s.companyRepo.Count(ctx, activeOnly)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Convertir a DTOs
	companyResponses := make([]*domainCompany.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		response := domainCompany.ToCompanyResponse(comp)
		
		// Obtener suscripción
		subscription, err := s.subscriptionRepo.FindByCompanyID(ctx, comp.ID)
		if err == nil && subscription != nil {
			response.Subscription = &domainCompany.SubscriptionInfo{
				ID:            subscription.ID,
				StartDate:     subscription.StartDate,
				EndDate:       subscription.EndDate,
				IsActive:      subscription.IsActive,
				DaysRemaining: subscription.DaysRemaining(),
			}
		}

		companyResponses = append(companyResponses, &response)
	}

	return &domainCompany.CompanyListResponse{
		Companies: companyResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// Activate activa una empresa
func (s *service) Activate(ctx context.Context, id uuid.UUID) error {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	if err := s.companyRepo.Activate(ctx, id); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// Deactivate desactiva una empresa
func (s *service) Deactivate(ctx context.Context, id uuid.UUID) error {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if company == nil {
		return errors.ErrNotFound.WithError(fmt.Errorf("company not found"))
	}

	if err := s.companyRepo.Deactivate(ctx, id); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}
