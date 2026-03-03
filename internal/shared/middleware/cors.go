package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORSConfig configuración CORS segura
type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig retorna configuración CORS por defecto
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Solo localhost en dev
		AllowCredentials: true,
		MaxAge:           3600, // 1 hora
	}
}

// ProductionCORSConfig retorna configuración para producción
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		MaxAge:           3600,
	}
}

// NewCORSMiddleware crea middleware CORS configurado
func NewCORSMiddleware(config CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: config.AllowCredentials,
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		MaxAge:           config.MaxAge,
	})
}

// SecurityHeaders agrega headers de seguridad
func SecurityHeaders() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Prevenir clickjacking
		c.Set("X-Frame-Options", "DENY")
		
		// Prevenir MIME sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// Habilitar XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy básica
		c.Set("Content-Security-Policy", "default-src 'self'")
		
		// HSTS (solo en producción con HTTPS)
		// c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		return c.Next()
	}
}
