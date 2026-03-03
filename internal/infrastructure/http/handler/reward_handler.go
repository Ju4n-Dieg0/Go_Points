package handler

import (
	"strconv"

	"github.com/Ju4n-Dieg0/Go_Points/internal/application/reward"
	rewardDomain "github.com/Ju4n-Dieg0/Go_Points/internal/domain/reward"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// RewardHandler maneja las peticiones HTTP de recompensas
type RewardHandler struct {
	service  *reward.RewardService
	validate *validator.Validate
}

// NewRewardHandler crea una nueva instancia
func NewRewardHandler(service *reward.RewardService) *RewardHandler {
	return &RewardHandler{
		service:  service,
		validate: validator.New(),
	}
}
// CreateReward crea una nueva recompensa
func (h *RewardHandler) CreateReward(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	var req rewardDomain.CreateRewardRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validar request
	if err := h.validate.Struct(req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Crear recompensa
	resp, err := h.service.CreateReward(c.Context(), companyID, &req)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdateReward actualiza una recompensa
func (h *RewardHandler) UpdateReward(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	var req rewardDomain.UpdateRewardRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validar request
	if err := h.validate.Struct(req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Actualizar recompensa
	resp, err := h.service.UpdateReward(c.Context(), id, companyID, &req)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteReward elimina una recompensa
func (h *RewardHandler) DeleteReward(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.service.DeleteReward(c.Context(), id, companyID); err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetReward obtiene una recompensa por ID
func (h *RewardHandler) GetReward(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	resp, err := h.service.GetReward(c.Context(), id, companyID)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ListRewards lista recompensas de una compañía
func (h *RewardHandler) ListRewards(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	resp, err := h.service.ListRewards(c.Context(), companyID, page, pageSize)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// CreateRewardPath crea un nuevo camino
func (h *RewardHandler) CreateRewardPath(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	var req rewardDomain.CreateRewardPathRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validar request
	if err := h.validate.Struct(req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Crear camino
	resp, err := h.service.CreateRewardPath(c.Context(), companyID, &req)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdateRewardPath actualiza un camino
func (h *RewardHandler) UpdateRewardPath(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	var req rewardDomain.UpdateRewardPathRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validar request
	if err := h.validate.Struct(req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Actualizar camino
	resp, err := h.service.UpdateRewardPath(c.Context(), id, companyID, &req)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteRewardPath elimina un camino
func (h *RewardHandler) DeleteRewardPath(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.service.DeleteRewardPath(c.Context(), id, companyID); err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetRewardPath obtiene un camino por ID
func (h *RewardHandler) GetRewardPath(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	resp, err := h.service.GetRewardPath(c.Context(), id, companyID)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ListRewardPaths lista caminos de una compañía
func (h *RewardHandler) ListRewardPaths(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	resp, err := h.service.ListRewardPaths(c.Context(), companyID, page, pageSize)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ReorderPathItems reordena los items de un camino
func (h *RewardHandler) ReorderPathItems(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	pathID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	var req rewardDomain.ReorderPathItemsRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validar request
	if err := h.validate.Struct(req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Reordenar items
	resp, err := h.service.ReorderPathItems(c.Context(), pathID, companyID, &req)
	if err != nil {
		return handleRewardError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// handleRewardError maneja errores de recompensas
func handleRewardError(c fiber.Ctx, err error) error {
	switch err {
	case reward.ErrRewardNotFound, reward.ErrPathNotFound:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	case reward.ErrInvalidCompany, reward.ErrInactiveSubscription:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	case reward.ErrProductNotFound, reward.ErrProductCompanyMismatch, reward.ErrRewardCompanyMismatch:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error interno del servidor",
		})
	}
}
