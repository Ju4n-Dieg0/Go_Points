package handler

import (
	"strconv"

	"github.com/Ju4n-Dieg0/Go_Points/internal/application/point"
	domainPoint "github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// PointHandler maneja las peticiones HTTP de puntos
type PointHandler struct {
	service  point.Service
	validate *validator.Validate
}

// NewPointHandler crea una nueva instancia del handler
func NewPointHandler(service point.Service) *PointHandler {
	return &PointHandler{
		service:  service,
		validate: validator.New(),
	}
}

// EarnPoints maneja la ganancia de puntos
func (h *PointHandler) EarnPoints(c fiber.Ctx) error {
	var req domainPoint.EarnPointsRequest

	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	response, err := h.service.EarnPoints(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// RedeemPoints maneja la redención de puntos
func (h *PointHandler) RedeemPoints(c fiber.Ctx) error {
	var req domainPoint.RedeemPointsRequest

	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	response, err := h.service.RedeemPoints(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetBalance obtiene el balance de puntos
func (h *PointHandler) GetBalance(c fiber.Ctx) error {
	consumerIDStr := c.Query("consumer_id")
	companyIDStr := c.Query("company_id")

	if consumerIDStr == "" || companyIDStr == "" {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "consumer_id and company_id are required"))
	}

	consumerID, err := uuid.Parse(consumerIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid consumer_id"))
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company_id"))
	}

	response, err := h.service.GetBalance(c.Context(), consumerID, companyID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetTransactions obtiene el historial de transacciones
func (h *PointHandler) GetTransactions(c fiber.Ctx) error {
	consumerIDStr := c.Query("consumer_id")
	companyIDStr := c.Query("company_id")

	if consumerIDStr == "" || companyIDStr == "" {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "consumer_id and company_id are required"))
	}

	consumerID, err := uuid.Parse(consumerIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid consumer_id"))
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company_id"))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.GetTransactions(c.Context(), consumerID, companyID, page, pageSize)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ConfigureRank configura los rangos de una empresa
func (h *PointHandler) ConfigureRank(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	var req domainPoint.ConfigureRankRequest

	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	response, err := h.service.ConfigureRank(c.Context(), companyID, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetRankConfig obtiene la configuración de rangos
func (h *PointHandler) GetRankConfig(c fiber.Ctx) error {
	companyIDStr := c.Params("companyId")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company_id"))
	}

	response, err := h.service.GetRankConfig(c.Context(), companyID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// handleValidationError convierte errores de validación
func (h *PointHandler) handleValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.ErrValidation.WithError(err)
	}

	details := make(map[string]string)
	for _, fieldErr := range validationErrors {
		details[fieldErr.Field()] = h.getValidationMessage(fieldErr)
	}

	return errors.ErrValidation.WithDetails(details)
}

// getValidationMessage retorna mensaje personalizado
func (h *PointHandler) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "Value must be greater than " + err.Param()
	default:
		return "Invalid value"
	}
}
