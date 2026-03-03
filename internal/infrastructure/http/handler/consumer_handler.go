package handler

import (
	"strconv"

	"github.com/Ju4n-Dieg0/Go_Points/internal/application/consumer"
	domainConsumer "github.com/Ju4n-Dieg0/Go_Points/internal/domain/consumer"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ConsumerHandler maneja las peticiones HTTP de consumidores
type ConsumerHandler struct {
	service  consumer.Service
	validate *validator.Validate
}

// NewConsumerHandler crea una nueva instancia del handler de consumidores
func NewConsumerHandler(service consumer.Service) *ConsumerHandler {
	return &ConsumerHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Create maneja la creación de un consumidor
func (h *ConsumerHandler) Create(c fiber.Ctx) error {
	var req domainConsumer.CreateConsumerRequest

	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	response, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetByID obtiene un consumidor por ID
func (h *ConsumerHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid consumer ID"))
	}

	response, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetByDocumentNumber obtiene un consumidor por número de documento
func (h *ConsumerHandler) GetByDocumentNumber(c fiber.Ctx) error {
	documentNumber := c.Params("documentNumber")
	if documentNumber == "" {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "document number is required"))
	}

	response, err := h.service.GetByDocumentNumber(c.Context(), documentNumber)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Update actualiza un consumidor
func (h *ConsumerHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid consumer ID"))
	}

	var req domainConsumer.UpdateConsumerRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	response, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Delete elimina un consumidor
func (h *ConsumerHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid consumer ID"))
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainConsumer.MessageResponse{
		Success: true,
		Message: "Consumer deleted successfully",
	})
}

// List obtiene una lista de consumidores con paginación
func (h *ConsumerHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.List(c.Context(), page, pageSize)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Search busca consumidores
func (h *ConsumerHandler) Search(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "search query is required"))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.Search(c.Context(), query, page, pageSize)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// handleValidationError convierte errores de validación en respuestas estructuradas
func (h *ConsumerHandler) handleValidationError(err error) error {
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

// getValidationMessage retorna un mensaje de error personalizado
func (h *ConsumerHandler) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum " + err.Param() + ")"
	case "url":
		return "Invalid URL format"
	case "oneof":
		return "Invalid value. Allowed values: " + err.Param()
	default:
		return "Invalid value"
	}
}
