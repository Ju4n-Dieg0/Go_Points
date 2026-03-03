package routes

import (
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupRewardRoutes configura las rutas de recompensas
func SetupRewardRoutes(api fiber.Router, rewardHandler *handler.RewardHandler, jwtConfig *config.JWTConfig) {
	// Rutas de recompensas (requieren rol COMPANY)
	rewards := api.Group("/rewards")
	rewards.Use(middleware.AuthMiddleware(jwtConfig))
	rewards.Use(middleware.RequireRole("COMPANY"))

	// CRUD de recompensas
	rewards.Post("/", rewardHandler.CreateReward)
	rewards.Get("/", rewardHandler.ListRewards)
	rewards.Get("/:id", rewardHandler.GetReward)
	rewards.Put("/:id", rewardHandler.UpdateReward)
	rewards.Delete("/:id", rewardHandler.DeleteReward)

	// CRUD de caminos
	paths := rewards.Group("/paths")
	paths.Post("/", rewardHandler.CreateRewardPath)
	paths.Get("/", rewardHandler.ListRewardPaths)
	paths.Get("/:id", rewardHandler.GetRewardPath)
	paths.Put("/:id", rewardHandler.UpdateRewardPath)
	paths.Delete("/:id", rewardHandler.DeleteRewardPath)

	// Reordenar items de un camino
	paths.Put("/:id/reorder", rewardHandler.ReorderPathItems)
}
