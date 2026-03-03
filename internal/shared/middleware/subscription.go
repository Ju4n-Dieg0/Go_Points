package middleware

import (
	subscriptionService "github.com/Ju4n-Dieg0/Go_Points/internal/application/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// RequireActiveSubscription middleware que valida que la empresa tenga suscripción activa
// Este middleware debe usarse en rutas que requieran suscripción activa (ej: registrar puntos, redimir puntos)
func RequireActiveSubscription(subService subscriptionService.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Obtener company_id del body, params o query
		var companyID uuid.UUID
		var err error

		// Intentar obtener de params
		companyIDStr := c.Params("companyId")
		if companyIDStr == "" {
			// Intentar obtener del body
			type requestBody struct {
				CompanyID string `json:"company_id"`
			}
			var body requestBody
			if err := c.Bind().JSON(&body); err == nil && body.CompanyID != "" {
				companyIDStr = body.CompanyID
			}
		}

		// Intentar obtener de query
		if companyIDStr == "" {
			companyIDStr = c.Query("company_id")
		}

		if companyIDStr == "" {
			return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "company_id is required"))
		}

		// Parsear UUID
		companyID, err = uuid.Parse(companyIDStr)
		if err != nil {
			return errors.ErrBadRequest.WithError(fiber.NewError(fiber.StatusBadRequest, "invalid company_id format"))
		}

		// Validar suscripción activa
		if err := subService.ValidateActiveSubscription(c.Context(), companyID); err != nil {
			return err
		}

		// Guardar company_id en el contexto para uso posterior
		c.Locals("companyID", companyID)

		return c.Next()
	}
}

// GetCompanyID obtiene el company_id del contexto de Fiber
func GetCompanyID(c fiber.Ctx) (uuid.UUID, error) {
	companyID, ok := c.Locals("companyID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "company_id not found in context")
	}
	return companyID, nil
}
