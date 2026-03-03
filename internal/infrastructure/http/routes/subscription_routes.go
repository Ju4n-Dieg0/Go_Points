package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupSubscriptionRoutes configura las rutas de suscripciones
func SetupSubscriptionRoutes(router fiber.Router, subscriptionHandler *handler.SubscriptionHandler, jwtConfig *config.JWTConfig) {
	subscriptions := router.Group("/subscriptions")

	// Rutas protegidas (requieren autenticación)
	subscriptions.Use(middleware.AuthMiddleware(jwtConfig))

	// Obtener suscripción de una empresa
	subscriptions.Get("/company/:companyId", subscriptionHandler.GetByCompanyID)

	// Solo super admins y empresas pueden renovar su suscripción
	subscriptions.Post("/company/:companyId/renew", 
		middleware.RequireRole(auth.RoleSuperAdmin, auth.RoleCompany), 
		subscriptionHandler.Renew,
	)

	// Solo super admins pueden cancelar suscripciones
	subscriptions.Post("/company/:companyId/cancel", 
		middleware.RequireRole(auth.RoleSuperAdmin), 
		subscriptionHandler.Cancel,
	)

	// Endpoint administrativo para verificar suscripciones expiradas
	subscriptions.Post("/check-expired", 
		middleware.RequireRole(auth.RoleSuperAdmin), 
		subscriptionHandler.CheckExpired,
	)
}
