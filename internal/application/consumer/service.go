package consumer

import (
	"context"
	"fmt"

	domainConsumer "github.com/Ju4n-Dieg0/Go_Points/internal/domain/consumer"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/persistence"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/google/uuid"
)

// Service define la lógica de negocio de consumidores
type Service interface {
	Create(ctx context.Context, req *domainConsumer.CreateConsumerRequest) (*domainConsumer.ConsumerResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainConsumer.ConsumerResponse, error)
	GetByDocumentNumber(ctx context.Context, documentNumber string) (*domainConsumer.ConsumerResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *domainConsumer.UpdateConsumerRequest) (*domainConsumer.ConsumerResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) (*domainConsumer.ConsumerListResponse, error)
	Search(ctx context.Context, query string, page, pageSize int) (*domainConsumer.ConsumerListResponse, error)
}

type service struct {
	repo domainConsumer.Repository
}

// NewService crea una nueva instancia del servicio de consumidores
func NewService(repo domainConsumer.Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create crea un nuevo consumidor
func (s *service) Create(ctx context.Context, req *domainConsumer.CreateConsumerRequest) (*domainConsumer.ConsumerResponse, error) {
	// Validar tipo de documento
	if !domainConsumer.IsValidDocumentType(req.DocumentType) {
		return nil, errors.ErrValidation.WithError(fmt.Errorf("invalid document type"))
	}

	// Verificar si ya existe un consumidor con ese documento
	exists, err := s.repo.ExistsByDocumentNumber(ctx, req.DocumentNumber)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if exists {
		return nil, errors.ErrConflict.WithError(fmt.Errorf("consumer with document number %s already exists", req.DocumentNumber))
	}

	// Crear consumidor
	consumer := &domainConsumer.Consumer{
		ID:             uuid.New(),
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Photo:          req.Photo,
	}

	if err := s.repo.Create(ctx, consumer); err != nil {
		// Manejar error de unique violation específicamente
		var uniqueErr *persistence.UniqueViolationError
		if e, ok := err.(*persistence.UniqueViolationError); ok {
			uniqueErr = e
			return nil, errors.ErrConflict.WithError(fmt.Errorf(uniqueErr.Message))
		}
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainConsumer.ToConsumerResponse(consumer)
	return &response, nil
}

// GetByID obtiene un consumidor por ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domainConsumer.ConsumerResponse, error) {
	consumer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if consumer == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("consumer not found"))
	}

	response := domainConsumer.ToConsumerResponse(consumer)
	return &response, nil
}

// GetByDocumentNumber obtiene un consumidor por número de documento
func (s *service) GetByDocumentNumber(ctx context.Context, documentNumber string) (*domainConsumer.ConsumerResponse, error) {
	consumer, err := s.repo.FindByDocumentNumber(ctx, documentNumber)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if consumer == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("consumer not found"))
	}

	response := domainConsumer.ToConsumerResponse(consumer)
	return &response, nil
}

// Update actualiza un consumidor existente
func (s *service) Update(ctx context.Context, id uuid.UUID, req *domainConsumer.UpdateConsumerRequest) (*domainConsumer.ConsumerResponse, error) {
	consumer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if consumer == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("consumer not found"))
	}

	// Actualizar campos
	if req.Name != nil {
		consumer.Name = *req.Name
	}
	if req.Email != nil {
		consumer.Email = *req.Email
	}
	if req.Phone != nil {
		consumer.Phone = req.Phone
	}
	if req.Photo != nil {
		consumer.Photo = req.Photo
	}

	if err := s.repo.Update(ctx, consumer); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainConsumer.ToConsumerResponse(consumer)
	return &response, nil
}

// Delete elimina un consumidor (soft delete)
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	consumer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if consumer == nil {
		return errors.ErrNotFound.WithError(fmt.Errorf("consumer not found"))
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// List obtiene una lista de consumidores con paginación
func (s *service) List(ctx context.Context, page, pageSize int) (*domainConsumer.ConsumerListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	consumers, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Convertir a DTOs
	consumerResponses := make([]*domainConsumer.ConsumerResponse, 0, len(consumers))
	for _, cons := range consumers {
		response := domainConsumer.ToConsumerResponse(cons)
		consumerResponses = append(consumerResponses, &response)
	}

	return &domainConsumer.ConsumerListResponse{
		Consumers: consumerResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// Search busca consumidores por nombre, email o documento
func (s *service) Search(ctx context.Context, query string, page, pageSize int) (*domainConsumer.ConsumerListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	consumers, err := s.repo.Search(ctx, query, pageSize, offset)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	total, err := s.repo.CountSearch(ctx, query)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Convertir a DTOs
	consumerResponses := make([]*domainConsumer.ConsumerResponse, 0, len(consumers))
	for _, cons := range consumers {
		response := domainConsumer.ToConsumerResponse(cons)
		consumerResponses = append(consumerResponses, &response)
	}

	return &domainConsumer.ConsumerListResponse{
		Consumers: consumerResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}
