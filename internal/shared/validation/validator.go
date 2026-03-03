package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	alphanumRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

// Validator wrapper para go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator crea una nueva instancia del validador con reglas personalizadas
func NewValidator() *Validator {
	v := validator.New()
	
	// Registrar validaciones personalizadas
	v.RegisterValidation("strong_password", validateStrongPassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("alphanumeric", validateAlphanumeric)
	v.RegisterValidation("no_sql_injection", validateNoSQLInjection)
	
	return &Validator{validate: v}
}

// Validate valida una estructura
func (v *Validator) Validate(data interface{}) error {
	return v.validate.Struct(data)
}

// ValidateVar valida una variable individual
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// GetValidationErrors convierte errores de validación a formato legible
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			errors[field] = getErrorMessage(e)
		}
	}
	
	return errors
}

// getErrorMessage retorna mensaje de error amigable
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	case "strong_password":
		return "Password must contain at least 8 characters, one uppercase, one lowercase, one number and one special character"
	case "phone":
		return "Invalid phone number format"
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	case "alphanumeric":
		return fmt.Sprintf("%s must contain only letters and numbers", e.Field())
	case "no_sql_injection":
		return fmt.Sprintf("%s contains invalid characters", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}

// validateStrongPassword valida contraseña fuerte
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	if len(password) < 8 {
		return false
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validatePhone valida número telefónico
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

// validateAlphanumeric valida solo alfanumérico
func validateAlphanumeric(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return alphanumRegex.MatchString(value)
}

// validateNoSQLInjection valida que no contenga caracteres de SQL injection
func validateNoSQLInjection(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	
	// Lista de patrones peligrosos
	dangerousPatterns := []string{
		"--", "/*", "*/", "xp_", "sp_", "exec", "execute",
		"drop", "truncate", "delete", "insert", "update",
		"union", "select", "script", "javascript", "alert",
		"onerror", "onload", "<script", "</script>",
	}
	
	for _, pattern := range dangerousPatterns {
		if strings.Contains(value, pattern) {
			return false
		}
	}
	
	return true
}

// ValidatePassword valida requisitos de contraseña
func ValidatePassword(password string, minLength int, requireSpecial bool) error {
	if len(password) < minLength {
		return fmt.Errorf("password must be at least %d characters", minLength)
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if requireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

// SanitizeInput sanitiza entrada de usuario
func SanitizeInput(input string) string {
	// Remover espacios al inicio y final
	input = strings.TrimSpace(input)
	
	// Remover caracteres de control
	var builder strings.Builder
	for _, char := range input {
		if !unicode.IsControl(char) {
			builder.WriteRune(char)
		}
	}
	
	return builder.String()
}

// ValidateUUID valida formato UUID
func ValidateUUID(id string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(id))
}
