package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupCompanyRoutes configura las rutas de empresas
func SetupCompanyRoutes(router fiber.Router, companyHandler *handler.CompanyHandler, jwtConfig *config.JWTConfig) {
	companies := router.Group("/companies")

	// Rutas protegidas (requieren autenticación)
	companies.Use(middleware.AuthMiddleware(jwtConfig))

	// Solo super admins pueden crear, actualizar y eliminar empresas
	companies.Post("/", middleware.RequireRole(auth.RoleSuperAdmin), companyHandler.Create)
	companies.Put("/:id", middleware.RequireRole(auth.RoleSuperAdmin), companyHandler.Update)
	companies.Delete("/:id", middleware.RequireRole(auth.RoleSuperAdmin), companyHandler.Delete)
	companies.Post("/:id/activate", middleware.RequireRole(auth.RoleSuperAdmin), companyHandler.Activate)
	companies.Post("/:id/deactivate", middleware.RequireRole(auth.RoleSuperAdmin), companyHandler.Deactivate)

	// Cualquier usuario autenticado puede ver empresas
	companies.Get("/", companyHandler.List)
	companies.Get("/:id", companyHandler.GetByID)
}
