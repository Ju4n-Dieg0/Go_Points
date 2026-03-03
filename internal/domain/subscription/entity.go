package subscription

import (
	"time"

	"github.com/google/uuid"
)

// Subscription representa la entidad de suscripción en el dominio
type Subscription struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_company_subscription_unique;index:idx_subscription_company_active,priority:1;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StartDate time.Time `gorm:"not null;index:idx_subscription_start"`
	EndDate   time.Time `gorm:"not null;index:idx_subscription_end;index:idx_subscription_company_active,priority:3"`
	IsActive  bool      `gorm:"default:true;not null;index:idx_subscription_active;index:idx_subscription_company_active,priority:2"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}

// TableName especifica el nombre de la tabla
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsExpired verifica si la suscripción ha expirado
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.EndDate)
}

// IsValid verifica si la suscripción es válida (activa y no expirada)
func (s *Subscription) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}

// Renew renueva la suscripción por 30 días adicionales
func (s *Subscription) Renew() {
	// Si la suscripción ya expiró, renovar desde hoy
	if s.IsExpired() {
		s.StartDate = time.Now()
		s.EndDate = time.Now().AddDate(0, 0, 30)
	} else {
		// Si aún está activa, extender desde la fecha de finalización actual
		s.EndDate = s.EndDate.AddDate(0, 0, 30)
	}
	s.IsActive = true
}

// Cancel cancela la suscripción (no la elimina, solo la marca como inactiva)
func (s *Subscription) Cancel() {
	s.IsActive = false
}

// DaysRemaining retorna los días restantes de la suscripción
func (s *Subscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	duration := time.Until(s.EndDate)
	return int(duration.Hours() / 24)
}
