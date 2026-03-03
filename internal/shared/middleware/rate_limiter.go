package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// RateLimiter implementa un rate limiter simple basado en IP
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests por ventana
	window   time.Duration // ventana de tiempo
}

// Visitor representa un visitante con su contador de requests
type Visitor struct {
	count      int
	lastReset  time.Time
	mu         sync.Mutex
}

// NewRateLimiter crea una nueva instancia del rate limiter
func NewRateLimiter(requestsPerWindow int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     requestsPerWindow,
		window:   window,
	}

	// Limpiar visitantes antiguos cada minuto
	go rl.cleanupVisitors()

	return rl
}

// getVisitor obtiene o crea un visitor para una IP
func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &Visitor{
			lastReset: time.Now(),
		}
		rl.visitors[ip] = v
	}

	return v
}

// isAllowed verifica si el visitor puede hacer una request
func (v *Visitor) isAllowed(rate int, window time.Duration) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	
	// Reset del contador si pasó la ventana de tiempo
	if now.Sub(v.lastReset) > window {
		v.count = 0
		v.lastReset = now
	}

	// Verificar si excedió el límite
	if v.count >= rate {
		return false
	}

	v.count++
	return true
}

// cleanupVisitors elimina visitantes inactivos
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			visitor.mu.Lock()
			if time.Since(visitor.lastReset) > rl.window*2 {
				delete(rl.visitors, ip)
			}
			visitor.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// Middleware retorna el middleware de Fiber para rate limiting
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		visitor := rl.getVisitor(ip)

		if !visitor.isAllowed(rl.rate, rl.window) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		}

		return c.Next()
	}
}

// RateLimitConfig configuración para diferentes endpoints
type RateLimitConfig struct {
	AuthRequests     int           // Requests permitidos para auth
	AuthWindow       time.Duration // Ventana de tiempo para auth
	GeneralRequests  int           // Requests para endpoints generales
	GeneralWindow    time.Duration // Ventana para endpoints generales
}

// DefaultRateLimitConfig retorna configuración por defecto
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		AuthRequests:    5,              // 5 intentos de login
		AuthWindow:      1 * time.Minute, // por minuto
		GeneralRequests: 100,            // 100 requests
		GeneralWindow:   1 * time.Minute, // por minuto
	}
}
