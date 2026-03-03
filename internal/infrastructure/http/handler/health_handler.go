package handler

import (
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/database"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/dto"
	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct {
	db *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck verifica que la aplicación esté ejecutándose
// @Summary Health check básico
// @Description Verifica que la aplicación esté ejecutándose
// @Tags Health
// @Produce json
// @Success 200 {object} dto.APIResponse{data=object{status=string,timestamp=string}}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c fiber.Ctx) error {
	return dto.Success(c, fiber.StatusOK, fiber.Map{
		"status":    "ok",
		"service":   "Go Points API",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessCheck verifica que la aplicación esté lista para recibir tráfico
// @Summary Readiness check
// @Description Verifica que la aplicación y sus dependencias (DB) estén listas
// @Tags Health
// @Produce json
// @Success 200 {object} dto.APIResponse{data=object{status=string,database=string,timestamp=string}}
// @Failure 503 {object} dto.APIResponse{error=dto.APIError}
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c fiber.Ctx) error {
	// Verificar conexión a la base de datos
	if err := h.db.HealthCheck(); err != nil {
		return dto.ServiceUnavailable(c, "Database is not ready")
	}

	return dto.Success(c, fiber.StatusOK, fiber.Map{
		"status":    "ready",
		"service":   "Go Points API",
		"database":  "healthy",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// LivenessCheck verifica que la aplicación no esté en deadlock
// @Summary Liveness check
// @Description Verifica que la aplicación esté viva (usado por Kubernetes)
// @Tags Health
// @Produce json
// @Success 200 {object} dto.APIResponse{data=object{status=string}}
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c fiber.Ctx) error {
	return dto.Success(c, fiber.StatusOK, fiber.Map{
		"status": "alive",
	})
}
