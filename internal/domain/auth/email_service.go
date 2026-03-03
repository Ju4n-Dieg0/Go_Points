package auth

import "context"

// EmailService define la interfaz para el servicio de email desacoplado
type EmailService interface {
	// SendWelcomeEmail envía un email de bienvenida al usuario
	SendWelcomeEmail(ctx context.Context, to, firstName string) error

	// SendPasswordResetEmail envía un email con el token de reset de contraseña
	SendPasswordResetEmail(ctx context.Context, to, firstName, resetToken string) error

	// SendPasswordChangedEmail envía un email notificando el cambio de contraseña
	SendPasswordChangedEmail(ctx context.Context, to, firstName string) error

	// SendEmailVerification envía un email de verificación
	SendEmailVerification(ctx context.Context, to, firstName, verificationToken string) error
}
