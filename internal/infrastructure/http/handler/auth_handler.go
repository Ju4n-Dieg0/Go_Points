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
// @Description Crea una nueva cuenta de usuario en el sistema. El email debe ser único.
// @Description Los roles disponibles son: SUPER_ADMIN, COMPANY (empresa), CONSUMER (consumidor).
// @Description La contraseña debe tener mínimo 8 caracteres.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domainAuth.RegisterRequest true "Datos de registro del usuario"
// @Success 201 {object} domainAuth.AuthResponse "Usuario registrado exitosamente con tokens"
// @Failure 400 {object} errors.ErrorResponse "Datos de entrada inválidos"
// @Failure 409 {object} errors.ErrorResponse "El email ya está registrado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
// @Description Autentica un usuario con email y contraseña.
// @Description Retorna access token (válido 15 min) y refresh token (válido 7 días).
// @Description El access token debe incluirse en el header Authorization como: Bearer {token}
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domainAuth.LoginRequest true "Credenciales de inicio de sesión"
// @Success 200 {object} domainAuth.AuthResponse "Autenticación exitosa con tokens"
// @Failure 400 {object} errors.ErrorResponse "Datos de entrada inválidos"
// @Failure 401 {object} errors.ErrorResponse "Credenciales incorrectas"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
// @Description Genera un nuevo access token usando un refresh token válido.
// @Description Útil cuando el access token ha expirado pero el refresh token aún es válido.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domainAuth.RefreshTokenRequest true "Refresh token válido"
// @Success 200 {object} domainAuth.RefreshTokenResponse "Nuevo access token generado"
// @Failure 400 {object} errors.ErrorResponse "Refresh token inválido o expirado"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
// @Description Invalida la sesión del usuario autenticado.
// @Description Requiere autenticación con Bearer token en el header Authorization.
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domainAuth.MessageResponse "Sesión cerrada exitosamente"
// @Failure 401 {object} errors.ErrorResponse "No autenticado o token inválido"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
// @Description Envía un email con enlace para recuperar la contraseña.
// @Description Por seguridad, siempre retorna éxito aunque el email no exista.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domainAuth.RequestPasswordResetRequest true "Email del usuario"
// @Success 200 {object} domainAuth.MessageResponse "Solicitud procesada (email enviado si existe)"
// @Failure 400 {object} errors.ErrorResponse "Email inválido"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
// @Description Establece una nueva contraseña usando el token recibido por email.
// @Description El token es de un solo uso y tiene tiempo de expiración limitado.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domainAuth.ConfirmPasswordResetRequest true "Token y nueva contraseña"
// @Success 200 {object} domainAuth.MessageResponse "Contraseña actualizada exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Token inválido, expirado o contraseña débil"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
// @Description Obtiene la información del perfil del usuario autenticado.
// @Description Requiere autenticación con Bearer token en el header Authorization.
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domainAuth.UserDTO "Información del usuario"
// @Failure 401 {object} errors.ErrorResponse "No autenticado o token inválido"
// @Failure 404 {object} errors.ErrorResponse "Usuario no encontrado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
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
