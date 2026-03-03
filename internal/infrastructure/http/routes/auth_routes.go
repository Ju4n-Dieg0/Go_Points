package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupAuthRoutes configura las rutas de autenticación
func SetupAuthRoutes(router fiber.Router, authHandler *handler.AuthHandler, jwtConfig *config.JWTConfig) {
	auth := router.Group("/auth")

	// Rutas públicas (sin autenticación)
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/password-reset/request", authHandler.RequestPasswordReset)
	auth.Post("/password-reset/confirm", authHandler.ConfirmPasswordReset)

	// Rutas protegidas (requieren autenticación)
	protected := auth.Use(middleware.AuthMiddleware(jwtConfig))
	protected.Post("/logout", authHandler.Logout)
	protected.Get("/profile", authHandler.GetProfile)
}
