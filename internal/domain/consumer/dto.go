package consumer

import (
	"time"

	"github.com/google/uuid"
)

// CreateConsumerRequest DTO para crear un consumidor
type CreateConsumerRequest struct {
	DocumentType   DocumentType `json:"document_type" validate:"required,oneof=DNI CE PASSPORT RUC OTHER"`
	DocumentNumber string       `json:"document_number" validate:"required,min=5,max=50"`
	Name           string       `json:"name" validate:"required,min=3,max=255"`
	Email          string       `json:"email" validate:"required,email,max=255"`
	Phone          *string      `json:"phone,omitempty" validate:"omitempty,min=7,max=20"`
	Photo          *string      `json:"photo,omitempty" validate:"omitempty,url"`
}

// UpdateConsumerRequest DTO para actualizar un consumidor
type UpdateConsumerRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Email *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=7,max=20"`
	Photo *string `json:"photo,omitempty" validate:"omitempty,url"`
}

// ConsumerResponse DTO para respuesta de consumidor
type ConsumerResponse struct {
	ID             uuid.UUID    `json:"id"`
	DocumentType   DocumentType `json:"document_type"`
	DocumentNumber string       `json:"document_number"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Phone          *string      `json:"phone,omitempty"`
	Photo          *string      `json:"photo,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// ConsumerListResponse DTO para lista de consumidores con paginación
type ConsumerListResponse struct {
	Consumers []*ConsumerResponse `json:"consumers"`
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	PageSize  int                 `json:"page_size"`
}

// SearchConsumerRequest DTO para búsqueda de consumidor
type SearchConsumerRequest struct {
	DocumentNumber string `json:"document_number" validate:"required,min=5,max=50"`
}

// ToConsumerResponse convierte Consumer entity a ConsumerResponse DTO
func ToConsumerResponse(consumer *Consumer) ConsumerResponse {
	return ConsumerResponse{
		ID:             consumer.ID,
		DocumentType:   consumer.DocumentType,
		DocumentNumber: consumer.DocumentNumber,
		Name:           consumer.Name,
		Email:          consumer.Email,
		Phone:          consumer.Phone,
		Photo:          consumer.Photo,
		CreatedAt:      consumer.CreatedAt,
		UpdatedAt:      consumer.UpdatedAt,
	}
}

// MessageResponse respuesta de mensaje simple
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
