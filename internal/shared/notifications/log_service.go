package notifications

import (
	"context"
	"log/slog"
)

// LogNotificationService implementación de notificaciones usando logger
type LogNotificationService struct {
	logger *slog.Logger
}

// NewLogNotificationService crea una nueva instancia del servicio de log
func NewLogNotificationService(logger *slog.Logger) *LogNotificationService {
	return &LogNotificationService{
		logger: logger,
	}
}

// NotifyPointExpiring registra notificación de puntos próximos a expirar
func (s *LogNotificationService) NotifyPointExpiring(ctx context.Context, data PointExpiringData) error {
	s.logger.InfoContext(ctx, "Point expiring notification",
		slog.String("notification_type", string(NotificationTypePointExpiring)),
		slog.String("consumer_id", data.ConsumerID.String()),
		slog.String("consumer_email", data.ConsumerEmail),
		slog.String("consumer_name", data.ConsumerName),
		slog.String("company_id", data.CompanyID.String()),
		slog.String("company_name", data.CompanyName),
		slog.Int64("points", data.Points),
		slog.Time("expiration_date", data.ExpirationDate),
		slog.Int("days_until_expiry", data.DaysUntilExpiry),
	)
	return nil
}

// NotifyRedemption registra notificación de redención de puntos
func (s *LogNotificationService) NotifyRedemption(ctx context.Context, data RedemptionData) error {
	s.logger.InfoContext(ctx, "Redemption notification",
		slog.String("notification_type", string(NotificationTypeRedemption)),
		slog.String("consumer_id", data.ConsumerID.String()),
		slog.String("consumer_email", data.ConsumerEmail),
		slog.String("consumer_name", data.ConsumerName),
		slog.String("company_id", data.CompanyID.String()),
		slog.String("company_name", data.CompanyName),
		slog.Int64("points_redeemed", data.PointsRedeemed),
		slog.Int64("remaining_balance", data.RemainingBalance),
		slog.String("transaction_id", data.TransactionID.String()),
		slog.Time("redeemed_at", data.RedeemedAt),
	)
	return nil
}

// NotifyRankUpgrade registra notificación de cambio de rango
func (s *LogNotificationService) NotifyRankUpgrade(ctx context.Context, data RankUpgradeData) error {
	s.logger.InfoContext(ctx, "Rank upgrade notification",
		slog.String("notification_type", string(NotificationTypeRankUpgrade)),
		slog.String("consumer_id", data.ConsumerID.String()),
		slog.String("consumer_email", data.ConsumerEmail),
		slog.String("consumer_name", data.ConsumerName),
		slog.String("company_id", data.CompanyID.String()),
		slog.String("company_name", data.CompanyName),
		slog.String("old_rank", data.OldRank),
		slog.String("new_rank", data.NewRank),
		slog.Int64("total_points", data.TotalPoints),
		slog.Time("upgraded_at", data.UpgradedAt),
	)
	return nil
}
