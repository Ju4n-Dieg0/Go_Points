package handler

import (
	"mime/multipart"
	"strconv"

	"github.com/Ju4n-Dieg0/Go_Points/internal/application/product"
	domainProduct "github.com/Ju4n-Dieg0/Go_Points/internal/domain/product"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ProductHandler maneja las peticiones HTTP de productos
type ProductHandler struct {
	service  product.Service
	validate *validator.Validate
}

// NewProductHandler crea una nueva instancia del handler de productos
func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Create maneja la creación de un producto
func (h *ProductHandler) Create(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	var req domainProduct.CreateProductRequest

	// Parsear form data
	if err := c.Bind().Body(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Obtener archivo de foto (opcional)
	fileHeader, err := c.FormFile("photo")
	var file multipart.File
	if err == nil && fileHeader != nil {
		file, err = fileHeader.Open()
		if err != nil {
			return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "failed to open photo file"))
		}
		defer file.Close()
	}

	response, err := h.service.Create(c.Context(), companyID, &req, file, fileHeader)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetByID obtiene un producto por ID
func (h *ProductHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid product ID"))
	}

	response, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Update actualiza un producto
func (h *ProductHandler) Update(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid product ID"))
	}

	var req domainProduct.UpdateProductRequest

	// Parsear form data
	if err := c.Bind().Body(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Obtener archivo de foto (opcional)
	fileHeader, err := c.FormFile("photo")
	var file multipart.File
	if err == nil && fileHeader != nil {
		file, err = fileHeader.Open()
		if err != nil {
			return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "failed to open photo file"))
		}
		defer file.Close()
	}

	response, err := h.service.Update(c.Context(), id, companyID, &req, file, fileHeader)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Delete elimina un producto
func (h *ProductHandler) Delete(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid product ID"))
	}

	if err := h.service.Delete(c.Context(), id, companyID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainProduct.MessageResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}

// List obtiene productos de la empresa con paginación
func (h *ProductHandler) List(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.List(c.Context(), companyID, page, pageSize)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ListAll obtiene todos los productos visibles (público)
func (h *ProductHandler) ListAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.ListAll(c.Context(), page, pageSize)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Search busca productos de la empresa
func (h *ProductHandler) Search(c fiber.Ctx) error {
	// Obtener companyID del usuario autenticado
	companyID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	query := c.Query("q")
	if query == "" {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "search query is required"))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.Search(c.Context(), companyID, query, page, pageSize)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// handleValidationError convierte errores de validación en respuestas estructuradas
func (h *ProductHandler) handleValidationError(err error) error {
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
func (h *ProductHandler) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short (minimum " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum " + err.Param() + ")"
	case "gt":
		return "Value must be greater than " + err.Param()
	default:
		return "Invalid value"
	}
}
