package notifications

import (
	"context"
	"fmt"

	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
)

// EmailNotificationService implementación de notificaciones por email
type EmailNotificationService struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
}

// NewEmailNotificationService crea una nueva instancia del servicio de email
func NewEmailNotificationService(cfg *config.EmailConfig) *EmailNotificationService {
	return &EmailNotificationService{
		smtpHost:     cfg.SMTPHost,
		smtpPort:     cfg.SMTPPort,
		smtpUser:     cfg.SMTPUser,
		smtpPassword: cfg.SMTPPassword,
		fromEmail:    cfg.FromEmail,
		fromName:     cfg.FromName,
	}
}

// NotifyPointExpiring envía notificación de puntos próximos a expirar
func (s *EmailNotificationService) NotifyPointExpiring(ctx context.Context, data PointExpiringData) error {
	subject := fmt.Sprintf("⏰ Tus puntos en %s están por expirar", data.CompanyName)
	
	body := fmt.Sprintf(`
		Hola %s,

		Te informamos que tienes %d puntos en %s que expirarán pronto.

		📅 Fecha de expiración: %s
		⏳ Días restantes: %d
		💎 Puntos a expirar: %d

		No pierdas tus puntos, ¡úsalos antes de que expiren!

		Saludos,
		Equipo de %s
	`,
		data.ConsumerName,
		data.Points,
		data.CompanyName,
		data.ExpirationDate.Format("02/01/2006"),
		data.DaysUntilExpiry,
		data.Points,
		data.CompanyName,
	)

	// TODO: Implementar envío real de email usando SMTP o servicio externo (SendGrid, AWS SES, etc.)
	// Por ahora solo registramos que se enviaría
	fmt.Printf("[EMAIL] To: %s | Subject: %s\n", data.ConsumerEmail, subject)
	fmt.Printf("[EMAIL] Body: %s\n", body)

	return nil
}

// NotifyRedemption envía notificación de redención de puntos
func (s *EmailNotificationService) NotifyRedemption(ctx context.Context, data RedemptionData) error {
	subject := fmt.Sprintf("✅ Redención exitosa en %s", data.CompanyName)
	
	body := fmt.Sprintf(`
		Hola %s,

		¡Tu redención ha sido procesada exitosamente!

		💰 Puntos redimidos: %d
		📊 Balance restante: %d
		🆔 ID de transacción: %s
		📅 Fecha: %s

		Gracias por ser parte de nuestro programa de lealtad.

		Saludos,
		Equipo de %s
	`,
		data.ConsumerName,
		data.PointsRedeemed,
		data.RemainingBalance,
		data.TransactionID,
		data.RedeemedAt.Format("02/01/2006 15:04:05"),
		data.CompanyName,
	)

	// TODO: Implementar envío real de email
	fmt.Printf("[EMAIL] To: %s | Subject: %s\n", data.ConsumerEmail, subject)
	fmt.Printf("[EMAIL] Body: %s\n", body)

	return nil
}

// NotifyRankUpgrade envía notificación de cambio de rango
func (s *EmailNotificationService) NotifyRankUpgrade(ctx context.Context, data RankUpgradeData) error {
	subject := fmt.Sprintf("🎉 ¡Felicitaciones! Has subido de rango en %s", data.CompanyName)
	
	body := fmt.Sprintf(`
		¡Hola %s!

		¡Tenemos excelentes noticias! Has alcanzado un nuevo nivel en nuestro programa de lealtad.

		⬆️ Nuevo Rango: %s (anteriormente: %s)
		💎 Total de puntos acumulados: %d
		📅 Fecha de actualización: %s

		Con tu nuevo rango, disfrutas de beneficios exclusivos y mejores recompensas.

		¡Sigue acumulando puntos para desbloquear más beneficios!

		Saludos,
		Equipo de %s
	`,
		data.ConsumerName,
		data.NewRank,
		data.OldRank,
		data.TotalPoints,
		data.UpgradedAt.Format("02/01/2006 15:04:05"),
		data.CompanyName,
	)

	// TODO: Implementar envío real de email
	fmt.Printf("[EMAIL] To: %s | Subject: %s\n", data.ConsumerEmail, subject)
	fmt.Printf("[EMAIL] Body: %s\n", body)

	return nil
}

// sendEmail método privado para enviar emails (implementación futura)
func (s *EmailNotificationService) sendEmail(to, subject, body string) error {
	// TODO: Implementar con biblioteca como gomail, go-mail, o cliente de servicio externo
	// Ejemplo con SMTP básico:
	// 
	// auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)
	// 
	// msg := []byte(fmt.Sprintf(
	//     "From: %s <%s>\r\n"+
	//     "To: %s\r\n"+
	//     "Subject: %s\r\n"+
	//     "Content-Type: text/html; charset=UTF-8\r\n"+
	//     "\r\n"+
	//     "%s\r\n",
	//     s.fromName, s.fromEmail, to, subject, body,
	// ))
	// 
	// addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	// return smtp.SendMail(addr, auth, s.fromEmail, []string{to}, msg)

	return nil
}
