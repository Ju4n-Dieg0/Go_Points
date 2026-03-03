package handler

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/application/auth"
	domainAuth "github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// AuthHandler maneja las peticiones HTTP de autenticación
type AuthHandler struct {
	service  auth.Service
	validate *validator.Validate
}

// NewAuthHandler crea una nueva instancia del handler de autenticación
func NewAuthHandler(service auth.Service) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Register maneja el registro de nuevos usuarios
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domainAuth.RegisterRequest true "Register Request"
// @Success 201 {object} domainAuth.AuthResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req domainAuth.RegisterRequest

	// Parse request body
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Call service
	response, err := h.service.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Login maneja el inicio de sesión
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domainAuth.LoginRequest true "Login Request"
// @Success 200 {object} domainAuth.AuthResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req domainAuth.LoginRequest

	// Parse request body
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Call service
	response, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// RefreshToken maneja la renovación del access token
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domainAuth.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} domainAuth.RefreshTokenResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req domainAuth.RefreshTokenRequest

	// Parse request body
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Call service
	response, err := h.service.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Logout maneja el cierre de sesión
// @Summary Logout user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domainAuth.MessageResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	// Call service
	if err := h.service.Logout(c.Context(), userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainAuth.MessageResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

// RequestPasswordReset maneja la solicitud de reset de contraseña
// @Summary Request password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domainAuth.RequestPasswordResetRequest true "Request Password Reset"
// @Success 200 {object} domainAuth.MessageResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /auth/password-reset/request [post]
func (h *AuthHandler) RequestPasswordReset(c fiber.Ctx) error {
	var req domainAuth.RequestPasswordResetRequest

	// Parse request body
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Call service
	if err := h.service.RequestPasswordReset(c.Context(), &req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainAuth.MessageResponse{
		Success: true,
		Message: "If the email exists, a password reset link has been sent",
	})
}

// ConfirmPasswordReset maneja la confirmación del reset de contraseña
// @Summary Confirm password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domainAuth.ConfirmPasswordResetRequest true "Confirm Password Reset"
// @Success 200 {object} domainAuth.MessageResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /auth/password-reset/confirm [post]
func (h *AuthHandler) ConfirmPasswordReset(c fiber.Ctx) error {
	var req domainAuth.ConfirmPasswordResetRequest

	// Parse request body
	if err := c.Bind().JSON(&req); err != nil {
		return errors.ErrBadRequest.WithError(err)
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return h.handleValidationError(err)
	}

	// Call service
	if err := h.service.ConfirmPasswordReset(c.Context(), &req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainAuth.MessageResponse{
		Success: true,
		Message: "Password has been reset successfully",
	})
}

// GetProfile obtiene el perfil del usuario autenticado
// @Summary Get user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domainAuth.UserDTO
// @Failure 401 {object} errors.ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return errors.ErrUnauthorized.WithError(err)
	}

	// Call service
	user, err := h.service.GetUserByID(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// handleValidationError convierte errores de validación en respuestas estructuradas
func (h *AuthHandler) handleValidationError(err error) error {
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

// getValidationMessage retorna un mensaje de error personalizado para cada tipo de validación
func (h *AuthHandler) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum " + err.Param() + ")"
	case "oneof":
		return "Invalid value. Allowed values: " + err.Param()
	default:
		return "Invalid value"
	}
}
