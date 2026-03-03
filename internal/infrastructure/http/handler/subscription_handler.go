package handler

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/application/subscription"
	domainSubscription "github.com/Ju4n-Dieg0/Go_Points/internal/domain/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// SubscriptionHandler maneja las peticiones HTTP de suscripciones
type SubscriptionHandler struct {
	service  subscription.Service
	validate *validator.Validate
}

// NewSubscriptionHandler crea una nueva instancia del handler de suscripciones
func NewSubscriptionHandler(service subscription.Service) *SubscriptionHandler {
	return &SubscriptionHandler{
		service:  service,
		validate: validator.New(),
	}
}

// GetByCompanyID obtiene la suscripción de una empresa
func (h *SubscriptionHandler) GetByCompanyID(c fiber.Ctx) error {
	companyIDStr := c.Params("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	response, err := h.service.GetByCompanyID(c.Context(), companyID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Renew renueva la suscripción de una empresa
func (h *SubscriptionHandler) Renew(c fiber.Ctx) error {
	companyIDStr := c.Params("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	response, err := h.service.Renew(c.Context(), companyID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Cancel cancela la suscripción de una empresa
func (h *SubscriptionHandler) Cancel(c fiber.Ctx) error {
	companyIDStr := c.Params("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company ID"))
	}

	if err := h.service.Cancel(c.Context(), companyID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainSubscription.MessageResponse{
		Success: true,
		Message: "Subscription cancelled successfully. It will expire on its end date.",
	})
}

// CheckExpired verifica y desactiva suscripciones expiradas (endpoint administrativo)
func (h *SubscriptionHandler) CheckExpired(c fiber.Ctx) error {
	if err := h.service.CheckAndDeactivateExpired(c.Context()); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domainSubscription.MessageResponse{
		Success: true,
		Message: "Expired subscriptions have been processed",
	})
}
