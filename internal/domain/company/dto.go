package company

import (
	"time"

	"github.com/google/uuid"
)

// CreateCompanyRequest DTO para crear una empresa
type CreateCompanyRequest struct {
	Name string  `json:"name" validate:"required,min=3,max=255"`
	Logo *string `json:"logo,omitempty" validate:"omitempty,url"`
}

// UpdateCompanyRequest DTO para actualizar una empresa
type UpdateCompanyRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Logo *string `json:"logo,omitempty" validate:"omitempty,url"`
}

// CompanyResponse DTO para respuesta de empresa
type CompanyResponse struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"name"`
	Logo      *string             `json:"logo,omitempty"`
	IsActive  bool                `json:"is_active"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Subscription *SubscriptionInfo `json:"subscription,omitempty"`
}

// SubscriptionInfo información básica de suscripción incluida en CompanyResponse
type SubscriptionInfo struct {
	ID            uuid.UUID `json:"id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsActive      bool      `json:"is_active"`
	DaysRemaining int       `json:"days_remaining"`
}

// CompanyListResponse DTO para lista de empresas con paginación
type CompanyListResponse struct {
	Companies []*CompanyResponse `json:"companies"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// ToCompanyResponse convierte Company entity a CompanyResponse DTO
func ToCompanyResponse(company *Company) CompanyResponse {
	return CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Logo:      company.Logo,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}

// MessageResponse respuesta de mensaje simple
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
