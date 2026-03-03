package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupConsumerRoutes configura las rutas de consumidores
func SetupConsumerRoutes(api fiber.Router, consumerHandler *handler.ConsumerHandler, jwtConfig *config.JWTConfig) {
	consumers := api.Group("/consumers")

	// Middleware de autenticación para todas las rutas de consumidores
	consumers.Use(middleware.AuthMiddleware(jwtConfig))

	// Rutas CRUD
	consumers.Post("/", consumerHandler.Create)
	consumers.Get("/", consumerHandler.List)
	consumers.Get("/search", consumerHandler.Search)
	consumers.Get("/:id", consumerHandler.GetByID)
	consumers.Get("/document/:documentNumber", consumerHandler.GetByDocumentNumber)
	consumers.Put("/:id", consumerHandler.Update)
	consumers.Delete("/:id", consumerHandler.Delete)
}
