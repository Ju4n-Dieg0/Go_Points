package service

import (
	"context"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
)

// stubEmailService implementación stub del servicio de email (para desarrollo)
type stubEmailService struct{}

// NewStubEmailService crea una nueva instancia del servicio stub de email
func NewStubEmailService() auth.EmailService {
	return &stubEmailService{}
}

// SendWelcomeEmail simula el envío de email de bienvenida
func (s *stubEmailService) SendWelcomeEmail(ctx context.Context, to, firstName string) error {
	logger.Info("Email sent (stub)",
		"type", "welcome",
		"to", to,
		"firstName", firstName,
	)
	return nil
}

// SendPasswordResetEmail simula el envío de email de reset de contraseña
func (s *stubEmailService) SendPasswordResetEmail(ctx context.Context, to, firstName, resetToken string) error {
	logger.Info("Email sent (stub)",
		"type", "password_reset",
		"to", to,
		"firstName", firstName,
		"resetToken", resetToken,
	)
	return nil
}

// SendPasswordChangedEmail simula el envío de email de confirmación de cambio de contraseña
func (s *stubEmailService) SendPasswordChangedEmail(ctx context.Context, to, firstName string) error {
	logger.Info("Email sent (stub)",
		"type", "password_changed",
		"to", to,
		"firstName", firstName,
	)
	return nil
}

// SendEmailVerification simula el envío de email de verificación
func (s *stubEmailService) SendEmailVerification(ctx context.Context, to, firstName, verificationToken string) error {
	logger.Info("Email sent (stub)",
		"type", "email_verification",
		"to", to,
		"firstName", firstName,
		"verificationToken", verificationToken,
	)
	return nil
}
