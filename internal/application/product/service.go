package product

import (
	"context"
	"fmt"
	"mime/multipart"

	subscriptionService "github.com/Ju4n-Dieg0/Go_Points/internal/application/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/files"
	domainProduct "github.com/Ju4n-Dieg0/Go_Points/internal/domain/product"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service define la interfaz para operaciones de negocio de productos
type Service interface {
	// Create crea un nuevo producto (requiere suscripción activa)
	Create(ctx context.Context, companyID uuid.UUID, req *domainProduct.CreateProductRequest, file multipart.File, fileHeader *multipart.FileHeader) (*domainProduct.ProductResponse, error)

	// GetByID obtiene un producto por ID
	GetByID(ctx context.Context, id uuid.UUID) (*domainProduct.ProductResponse, error)

	// Update actualiza un producto existente
	Update(ctx context.Context, id, companyID uuid.UUID, req *domainProduct.UpdateProductRequest, file multipart.File, fileHeader *multipart.FileHeader) (*domainProduct.ProductResponse, error)

	// Delete elimina un producto
	Delete(ctx context.Context, id, companyID uuid.UUID) error

	// List lista productos de una empresa con paginación
	List(ctx context.Context, companyID uuid.UUID, page, pageSize int) (*domainProduct.ProductListResponse, error)

	// ListAll lista todos los productos visibles (para consumidores)
	ListAll(ctx context.Context, page, pageSize int) (*domainProduct.ProductListResponse, error)

	// Search busca productos de una empresa
	Search(ctx context.Context, companyID uuid.UUID, query string, page, pageSize int) (*domainProduct.ProductListResponse, error)
}

// service implementa Service
type service struct {
	repo               domainProduct.Repository
	fileService        files.FileService
	subscriptionService subscriptionService.Service
}

// NewService crea una nueva instancia de Service
func NewService(repo domainProduct.Repository, fileService files.FileService, subscriptionSvc subscriptionService.Service) Service {
	return &service{
		repo:               repo,
		fileService:        fileService,
		subscriptionService: subscriptionSvc,
	}
}

// Create crea un nuevo producto
func (s *service) Create(ctx context.Context, companyID uuid.UUID, req *domainProduct.CreateProductRequest, file multipart.File, fileHeader *multipart.FileHeader) (*domainProduct.ProductResponse, error) {
	// Validar que la empresa tiene suscripción activa
	if err := s.subscriptionService.ValidateActiveSubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Crear producto
	product := &domainProduct.Product{
		ID:          uuid.New(),
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsVisible:   true, // Por defecto visible
	}

	// Si se proporciona is_visible en el request, usarlo
	if req.IsVisible != nil {
		product.IsVisible = *req.IsVisible
	}

	// Subir foto si se proporciona
	if file != nil && fileHeader != nil {
		photoPath, err := s.fileService.Upload(ctx, file, fileHeader)
		if err != nil {
			return nil, err
		}
		product.Photo = photoPath
	}

	// Guardar en base de datos
	if err := s.repo.Create(ctx, product); err != nil {
		// Si hay error y se subió foto, intentar eliminarla
		if product.Photo != "" {
			_ = s.fileService.Delete(ctx, product.Photo)
		}
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainProduct.ToProductResponse(product)
	return &response, nil
}

// GetByID obtiene un producto por ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domainProduct.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithError(fmt.Errorf("product not found"))
		}
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainProduct.ToProductResponse(product)
	return &response, nil
}

// Update actualiza un producto existente
func (s *service) Update(ctx context.Context, id, companyID uuid.UUID, req *domainProduct.UpdateProductRequest, file multipart.File, fileHeader *multipart.FileHeader) (*domainProduct.ProductResponse, error) {
	// Validar que la empresa tiene suscripción activa
	if err := s.subscriptionService.ValidateActiveSubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Buscar producto
	product, err := s.repo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithError(fmt.Errorf("product not found"))
		}
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Actualizar campos si se proporcionan
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.IsVisible != nil {
		product.IsVisible = *req.IsVisible
	}

	// Subir nueva foto si se proporciona
	if file != nil && fileHeader != nil {
		// Eliminar foto anterior
		if product.Photo != "" {
			_ = s.fileService.Delete(ctx, product.Photo)
		}

		// Subir nueva foto
		photoPath, err := s.fileService.Upload(ctx, file, fileHeader)
		if err != nil {
			return nil, err
		}
		product.Photo = photoPath
	}

	// Guardar cambios
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainProduct.ToProductResponse(product)
	return &response, nil
}

// Delete elimina un producto
func (s *service) Delete(ctx context.Context, id, companyID uuid.UUID) error {
	// Validar que la empresa tiene suscripción activa
	if err := s.subscriptionService.ValidateActiveSubscription(ctx, companyID); err != nil {
		return err
	}

	// Buscar producto para obtener la foto
	product, err := s.repo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound.WithError(fmt.Errorf("product not found"))
		}
		return errors.ErrDatabase.WithError(err)
	}

	// Eliminar producto (soft delete)
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	// Eliminar foto si existe
	if product.Photo != "" {
		_ = s.fileService.Delete(ctx, product.Photo)
	}

	return nil
}

// List lista productos de una empresa con paginación
func (s *service) List(ctx context.Context, companyID uuid.UUID, page, pageSize int) (*domainProduct.ProductListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.repo.List(ctx, companyID, page, pageSize)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainProduct.ToProductListResponse(products, total, page, pageSize)
	return &response, nil
}

// ListAll lista todos los productos visibles (para consumidores)
func (s *service) ListAll(ctx context.Context, page, pageSize int) (*domainProduct.ProductListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.repo.ListAll(ctx, page, pageSize)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainProduct.ToProductListResponse(products, total, page, pageSize)
	return &response, nil
}

// Search busca productos de una empresa
func (s *service) Search(ctx context.Context, companyID uuid.UUID, query string, page, pageSize int) (*domainProduct.ProductListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.repo.Search(ctx, companyID, query, page, pageSize)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainProduct.ToProductListResponse(products, total, page, pageSize)
	return &response, nil
}
