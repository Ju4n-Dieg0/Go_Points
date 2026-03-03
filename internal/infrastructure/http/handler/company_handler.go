package handler

import (
	"strconv"

	"github.com/Ju4n-Dieg0/Go_Points/internal/application/company"
	domainCompany "github.com/Ju4n-Dieg0/Go_Points/internal/domain/company"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// CompanyHandler maneja las peticiones HTTP de empresas
type CompanyHandler struct {
	service  company.Service
	validate *validator.Validate
}

// NewCompanyHandler crea una nueva instancia del handler de empresas
func NewCompanyHandler(service company.Service) *CompanyHandler {
	return &CompanyHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Create maneja la creación de una empresa
func (h *CompanyHandler) Create(c fiber.Ctx) error {
	var req domainCompany.CreateCompanyRequest

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

// GetByID obtiene una empresa por ID
func (h *CompanyHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	response, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Update actualiza una empresa
func (h *CompanyHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	var req domainCompany.UpdateCompanyRequest
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

// Delete elimina una empresa
func (h *CompanyHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainCompany.MessageResponse{
		Success: true,
		Message: "Company deleted successfully",
	})
}

// List obtiene una lista de empresas con paginación
func (h *CompanyHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	activeOnly := c.Query("active_only", "false") == "true"

	response, err := h.service.List(c.Context(), page, pageSize, activeOnly)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Activate activa una empresa
func (h *CompanyHandler) Activate(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	if err := h.service.Activate(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainCompany.MessageResponse{
		Success: true,
		Message: "Company activated successfully",
	})
}

// Deactivate desactiva una empresa
func (h *CompanyHandler) Deactivate(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	if err := h.service.Deactivate(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainCompany.MessageResponse{
		Success: true,
		Message: "Company deactivated successfully",
	})
}

// handleValidationError convierte errores de validación en respuestas estructuradas
func (h *CompanyHandler) handleValidationError(err error) error {
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
func (h *CompanyHandler) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short (minimum " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum " + err.Param() + ")"
	case "url":
		return "Invalid URL format"
	default:
		return "Invalid value"
	}
}
