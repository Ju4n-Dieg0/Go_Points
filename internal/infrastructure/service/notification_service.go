package service

import (
	"context"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"github.com/google/uuid"
)

// StubNotificationService implementación stub de NotificationService
type StubNotificationService struct{}

// NewStubNotificationService crea una nueva instancia
func NewStubNotificationService() point.NotificationService {
	return &StubNotificationService{}
}

// NotifyPointsExpiring notifica que puntos están por expirar
func (s *StubNotificationService) NotifyPointsExpiring(ctx context.Context, consumerID, companyID uuid.UUID, points int64, expirationDate time.Time) error {
	logger.Info("Points expiring soon",
		"consumer_id", consumerID,
		"company_id", companyID,
		"points", points,
		"expiration_date", expirationDate,
	)
	// TODO: Implementar envío de email/SMS/push notification
	return nil
}

// NotifyPointsExpired notifica que puntos han expirado
func (s *StubNotificationService) NotifyPointsExpired(ctx context.Context, consumerID, companyID uuid.UUID, points int64) error {
	logger.Info("Points expired",
		"consumer_id", consumerID,
		"company_id", companyID,
		"points", points,
	)
	// TODO: Implementar envío de email/SMS/push notification
	return nil
}

// NotifyInactivityPenalty notifica penalización por inactividad
func (s *StubNotificationService) NotifyInactivityPenalty(ctx context.Context, consumerID, companyID uuid.UUID, penalty int64) error {
	logger.Info("Inactivity penalty applied",
		"consumer_id", consumerID,
		"company_id", companyID,
		"penalty", penalty,
	)
	// TODO: Implementar envío de email/SMS/push notification
	return nil
}
