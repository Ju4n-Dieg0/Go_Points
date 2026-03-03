package subscription

import (
	"time"

	"github.com/google/uuid"
)

// RenewSubscriptionRequest DTO para renovar suscripción
type RenewSubscriptionRequest struct {
	CompanyID uuid.UUID `json:"company_id" validate:"required,uuid"`
}

// CancelSubscriptionRequest DTO para cancelar suscripción
type CancelSubscriptionRequest struct {
	CompanyID uuid.UUID `json:"company_id" validate:"required,uuid"`
}

// SubscriptionResponse DTO para respuesta de suscripción
type SubscriptionResponse struct {
	ID            uuid.UUID `json:"id"`
	CompanyID     uuid.UUID `json:"company_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsActive      bool      `json:"is_active"`
	IsExpired     bool      `json:"is_expired"`
	DaysRemaining int       `json:"days_remaining"`
	CreatedAt     time.Time `json:"created_at"`
}

// SubscriptionListResponse DTO para lista de suscripciones
type SubscriptionListResponse struct {
	Subscriptions []*SubscriptionResponse `json:"subscriptions"`
	Total         int64                   `json:"total"`
}

// ToSubscriptionResponse convierte Subscription entity a SubscriptionResponse DTO
func ToSubscriptionResponse(sub *Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:            sub.ID,
		CompanyID:     sub.CompanyID,
		StartDate:     sub.StartDate,
		EndDate:       sub.EndDate,
		IsActive:      sub.IsActive,
		IsExpired:     sub.IsExpired(),
		DaysRemaining: sub.DaysRemaining(),
		CreatedAt:     sub.CreatedAt,
	}
}

// MessageResponse respuesta de mensaje simple
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
