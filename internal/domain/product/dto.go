package product

import (
	"time"

	"github.com/google/uuid"
)

// CreateProductRequest representa la solicitud de creación de producto
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	IsVisible   *bool   `json:"is_visible"`
}

// UpdateProductRequest representa la solicitud de actualización de producto
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
	IsVisible   *bool   `json:"is_visible"`
}

// ProductResponse representa la respuesta con datos de producto
type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Photo       string    `json:"photo,omitempty"`
	IsVisible   bool      `json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductListResponse representa la respuesta de listado de productos con paginación
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// MessageResponse representa una respuesta simple con mensaje
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ToProductResponse convierte una entidad Product a ProductResponse
func ToProductResponse(product *Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		CompanyID:   product.CompanyID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Photo:       product.Photo,
		IsVisible:   product.IsVisible,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToProductListResponse convierte una lista de productos a ProductListResponse
func ToProductListResponse(products []Product, total int64, page, pageSize int) ProductListResponse {
	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = ToProductResponse(&product)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return ProductListResponse{
		Products:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
