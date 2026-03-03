package point

import (
	"time"

	"github.com/google/uuid"
)

// EarnPointsRequest solicitud para ganar puntos
type EarnPointsRequest struct {
	ConsumerID uuid.UUID `json:"consumer_id" validate:"required"`
	CompanyID  uuid.UUID `json:"company_id" validate:"required"`
	Points     int64     `json:"points" validate:"required,gt=0"`
}

// RedeemPointsRequest solicitud para redimir puntos
type RedeemPointsRequest struct {
	ConsumerID uuid.UUID `json:"consumer_id" validate:"required"`
	CompanyID  uuid.UUID `json:"company_id" validate:"required"`
	Points     int64     `json:"points" validate:"required,gt=0"`
}

// ConfigureRankRequest solicitud para configurar rangos de empresa
type ConfigureRankRequest struct {
	SilverMinPoints int64 `json:"silver_min_points" validate:"required,gt=0"`
	GoldMinPoints   int64 `json:"gold_min_points" validate:"required,gt=0"`
}

// PointBalanceResponse respuesta con el balance de puntos
type PointBalanceResponse struct {
	ConsumerID            uuid.UUID  `json:"consumer_id"`
	CompanyID             uuid.UUID  `json:"company_id"`
	TotalHistoricalPoints int64      `json:"total_historical_points"`
	TotalAvailablePoints  int64      `json:"total_available_points"`
	Rank                  Rank       `json:"rank"`
	LastRedemptionDate    *time.Time `json:"last_redemption_date,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// TransactionResponse respuesta de una transacción
type TransactionResponse struct {
	ID              uuid.UUID       `json:"id"`
	ConsumerID      uuid.UUID       `json:"consumer_id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	Points          int64           `json:"points"`
	RemainingPoints int64           `json:"remaining_points"`
	Type            TransactionType `json:"type"`
	ExpirationDate  *time.Time      `json:"expiration_date,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// TransactionListResponse respuesta con lista de transacciones
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}

// RankConfigResponse respuesta con configuración de rangos
type RankConfigResponse struct {
	ID              uuid.UUID `json:"id"`
	CompanyID       uuid.UUID `json:"company_id"`
	SilverMinPoints int64     `json:"silver_min_points"`
	GoldMinPoints   int64     `json:"gold_min_points"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MessageResponse respuesta simple con mensaje
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ToPointBalanceResponse convierte ConsumerCompanyPoints a DTO
func ToPointBalanceResponse(balance *ConsumerCompanyPoints, rank Rank) PointBalanceResponse {
	return PointBalanceResponse{
		ConsumerID:            balance.ConsumerID,
		CompanyID:             balance.CompanyID,
		TotalHistoricalPoints: balance.TotalHistoricalPoints,
		TotalAvailablePoints:  balance.TotalAvailablePoints,
		Rank:                  rank,
		LastRedemptionDate:    balance.LastRedemptionDate,
		CreatedAt:             balance.CreatedAt,
		UpdatedAt:             balance.UpdatedAt,
	}
}

// ToTransactionResponse convierte PointTransaction a DTO
func ToTransactionResponse(tx *PointTransaction) TransactionResponse {
	return TransactionResponse{
		ID:              tx.ID,
		ConsumerID:      tx.ConsumerID,
		CompanyID:       tx.CompanyID,
		Points:          tx.Points,
		RemainingPoints: tx.RemainingPoints,
		Type:            tx.Type,
		ExpirationDate:  tx.ExpirationDate,
		CreatedAt:       tx.CreatedAt,
	}
}

// ToTransactionListResponse convierte lista de transacciones a DTO
func ToTransactionListResponse(transactions []PointTransaction, total int64, page, pageSize int) TransactionListResponse {
	responses := make([]TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = ToTransactionResponse(&tx)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return TransactionListResponse{
		Transactions: responses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}
}

// ToRankConfigResponse convierte CompanyRankConfig a DTO
func ToRankConfigResponse(config *CompanyRankConfig) RankConfigResponse {
	return RankConfigResponse{
		ID:              config.ID,
		CompanyID:       config.CompanyID,
		SilverMinPoints: config.SilverMinPoints,
		GoldMinPoints:   config.GoldMinPoints,
		CreatedAt:       config.CreatedAt,
		UpdatedAt:       config.UpdatedAt,
	}
}
