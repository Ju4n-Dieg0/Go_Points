package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupPointRoutes configura las rutas de puntos
func SetupPointRoutes(api fiber.Router, pointHandler *handler.PointHandler, jwtConfig *config.JWTConfig) {
	points := api.Group("/points")

	// Middleware de autenticación para todas las rutas
	points.Use(middleware.AuthMiddleware(jwtConfig))

	// Rutas de puntos
	points.Post("/earn", pointHandler.EarnPoints)           // Ganar puntos (COMPANY)
	points.Post("/redeem", pointHandler.RedeemPoints)       // Redimir puntos (CONSUMER/COMPANY)
	points.Get("/balance", pointHandler.GetBalance)         // Ver balance
	points.Get("/transactions", pointHandler.GetTransactions) // Historial

	// Configuración de rangos (solo COMPANY)
	companyOnly := points.Group("/")
	companyOnly.Use(middleware.RequireRole("COMPANY"))
	companyOnly.Post("/rank/configure", pointHandler.ConfigureRank)
	companyOnly.Get("/rank/:companyId", pointHandler.GetRankConfig)
}
