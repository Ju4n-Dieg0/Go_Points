package reward

import (
	"time"

	"github.com/google/uuid"
)

// CreateRewardRequest solicitud para crear recompensa
type CreateRewardRequest struct {
	ProductID      uuid.UUID `json:"product_id" validate:"required"`
	RequiredPoints int64     `json:"required_points" validate:"required,gt=0"`
}

// UpdateRewardRequest solicitud para actualizar recompensa
type UpdateRewardRequest struct {
	ProductID      *uuid.UUID `json:"product_id" validate:"omitempty"`
	RequiredPoints *int64     `json:"required_points" validate:"omitempty,gt=0"`
}

// CreateRewardPathRequest solicitud para crear camino de recompensas
type CreateRewardPathRequest struct {
	Name      string      `json:"name" validate:"required,min=3,max=255"`
	RewardIDs []uuid.UUID `json:"reward_ids" validate:"required,min=1"`
}

// UpdateRewardPathRequest solicitud para actualizar camino
type UpdateRewardPathRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=255"`
}
// ReorderPathItemsRequest solicitud para reordenar items del camino
type ReorderPathItemsRequest struct {
	Items []PathItemOrder `json:"items" validate:"required,min=1"`
}

// PathItemOrder define el orden de un item
type PathItemOrder struct {
	RewardID uuid.UUID `json:"reward_id" validate:"required"`
	Order    int       `json:"order" validate:"required,gte=0"`
}

// RewardResponse respuesta de recompensa
type RewardResponse struct {
	ID             uuid.UUID `json:"id"`
	CompanyID      uuid.UUID `json:"company_id"`
	ProductID      uuid.UUID `json:"product_id"`
	ProductName    string    `json:"product_name,omitempty"`
	RequiredPoints int64     `json:"required_points"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RewardListResponse respuesta con lista de recompensas
type RewardListResponse struct {
	Rewards    []RewardResponse `json:"rewards"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// RewardPathResponse respuesta de camino de recompensas
type RewardPathResponse struct {
	ID        uuid.UUID        `json:"id"`
	CompanyID uuid.UUID        `json:"company_id"`
	Name      string           `json:"name"`
	Items     []PathItemDetail `json:"items,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// PathItemDetail detalle de un item del camino
type PathItemDetail struct {
	ID             uuid.UUID `json:"id"`
	RewardID       uuid.UUID `json:"reward_id"`
	ProductID      uuid.UUID `json:"product_id"`
	ProductName    string    `json:"product_name"`
	RequiredPoints int64     `json:"required_points"`
	Order          int       `json:"order"`
}

// RewardPathListResponse respuesta con lista de caminos
type RewardPathListResponse struct {
	Paths      []RewardPathResponse `json:"paths"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// MessageResponse respuesta simple
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ToRewardResponse convierte Reward a DTO
func ToRewardResponse(reward *Reward) *RewardResponse {
	return &RewardResponse{
		ID:             reward.ID,
		CompanyID:      reward.CompanyID,
		ProductID:      reward.ProductID,
		ProductName:    "", // TODO: Fetch product name if needed
		RequiredPoints: reward.RequiredPoints,
		CreatedAt:      reward.CreatedAt,
		UpdatedAt:      reward.UpdatedAt,
	}
}

// ToRewardListResponse convierte lista de rewards a DTO
func ToRewardListResponse(rewards []Reward, total int64, page, pageSize int) *RewardListResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	rewardResponses := make([]RewardResponse, len(rewards))
	for i, r := range rewards {
		rewardResponses[i] = *ToRewardResponse(&r)
	}

	return &RewardListResponse{
		Rewards:    rewardResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToRewardPathResponse convierte RewardPath a DTO
func ToRewardPathResponse(path *RewardPath, items []PathItemDetail) *RewardPathResponse {
	return &RewardPathResponse{
		ID:        path.ID,
		CompanyID: path.CompanyID,
		Name:      path.Name,
		Items:     items,
		CreatedAt: path.CreatedAt,
		UpdatedAt: path.UpdatedAt,
	}
}

// ToRewardPathListResponse convierte lista de paths a DTO
func ToRewardPathListResponse(paths []RewardPathResponse, total int64, page, pageSize int) *RewardPathListResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &RewardPathListResponse{
		Paths:      paths,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}