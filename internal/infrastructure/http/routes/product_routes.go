package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupProductRoutes configura las rutas de productos
func SetupProductRoutes(api fiber.Router, productHandler *handler.ProductHandler, jwtConfig *config.JWTConfig) {
	products := api.Group("/products")

	// Rutas públicas (sin autenticación)
	products.Get("/catalog", productHandler.ListAll) // Catálogo público de productos visibles

	// Rutas protegidas (requieren autenticación y rol COMPANY)
	protected := products.Use(middleware.AuthMiddleware(jwtConfig))
	protected.Use(middleware.RequireRole("COMPANY")) // Solo empresas pueden gestionar productos

	protected.Post("/", productHandler.Create)
	protected.Get("/", productHandler.List)
	protected.Get("/search", productHandler.Search)
	protected.Get("/:id", productHandler.GetByID)
	protected.Put("/:id", productHandler.Update)
	protected.Delete("/:id", productHandler.Delete)
}
