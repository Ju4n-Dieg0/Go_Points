package middleware

import (
	"strings"

	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims representa los claims del JWT
type JWTClaims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Role   auth.Role `json:"role"`
	Type   string    `json:"type"`
	jwt.RegisteredClaims
}

// AuthMiddleware middleware de autenticación JWT para Fiber v3
func AuthMiddleware(jwtConfig *config.JWTConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extraer token del header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.ErrUnauthorized.WithError(fiber.NewError(fiber.StatusUnauthorized, "missing authorization header"))
		}

		// Verificar formato Bearer
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.ErrUnauthorized.WithError(fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format"))
		}

		tokenString := parts[1]

		// Parsear y validar token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verificar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
			}
			return []byte(jwtConfig.AccessSecret), nil
		})

		if err != nil {
			return errors.ErrUnauthorized.WithError(fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token"))
		}

		// Extraer claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return errors.ErrUnauthorized.WithError(fiber.NewError(fiber.StatusUnauthorized, "invalid token claims"))
		}

		// Verificar que sea un access token
		if claims.Type != "access" {
			return errors.ErrUnauthorized.WithError(fiber.NewError(fiber.StatusUnauthorized, "invalid token type"))
		}

		// Guardar información del usuario en el contexto
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// RequireRole middleware que requiere roles específicos
func RequireRole(roles ...auth.Role) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(auth.Role)
		if !ok {
			return errors.ErrUnauthorized.WithError(fiber.NewError(fiber.StatusUnauthorized, "user role not found"))
		}

		// Verificar si el usuario tiene uno de los roles requeridos
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return errors.ErrForbidden.WithError(fiber.NewError(fiber.StatusForbidden, "insufficient permissions"))
		}

		return c.Next()
	}
}

// RequireSuperAdmin middleware que requiere rol de super admin
func RequireSuperAdmin() fiber.Handler {
	return RequireRole(auth.RoleSuperAdmin)
}

// RequireCompany middleware que requiere rol de empresa
func RequireCompany() fiber.Handler {
	return RequireRole(auth.RoleCompany, auth.RoleSuperAdmin)
}

// RequireConsumer middleware que requiere rol de consumidor
func RequireConsumer() fiber.Handler {
	return RequireRole(auth.RoleConsumer, auth.RoleSuperAdmin)
}

// GetUserID obtiene el ID del usuario del contexto
func GetUserID(c fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "user ID not found in context")
	}
	return userID, nil
}

// GetUserEmail obtiene el email del usuario del contexto
func GetUserEmail(c fiber.Ctx) (string, error) {
	email, ok := c.Locals("userEmail").(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user email not found in context")
	}
	return email, nil
}

// GetUserRole obtiene el rol del usuario del contexto
func GetUserRole(c fiber.Ctx) (auth.Role, error) {
	role, ok := c.Locals("userRole").(auth.Role)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user role not found in context")
	}
	return role, nil
}
