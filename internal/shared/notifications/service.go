package notifications

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NotificationType representa el tipo de notificación
type NotificationType string

const (
	NotificationTypePointExpiring NotificationType = "POINT_EXPIRING"
	NotificationTypeRedemption    NotificationType = "REDEMPTION"
	NotificationTypeRankUpgrade   NotificationType = "RANK_UPGRADE"
)

// PointExpiringData datos para notificación de puntos a expirar
type PointExpiringData struct {
	ConsumerID      uuid.UUID
	ConsumerEmail   string
	ConsumerName    string
	CompanyID       uuid.UUID
	CompanyName     string
	Points          int64
	ExpirationDate  time.Time
	DaysUntilExpiry int
}

// RedemptionData datos para notificación de redención
type RedemptionData struct {
	ConsumerID    uuid.UUID
	ConsumerEmail string
	ConsumerName  string
	CompanyID     uuid.UUID
	CompanyName   string
	PointsRedeemed int64
	RemainingBalance int64
	TransactionID uuid.UUID
	RedeemedAt    time.Time
}

// RankUpgradeData datos para notificación de cambio de rango
type RankUpgradeData struct {
	ConsumerID    uuid.UUID
	ConsumerEmail string
	ConsumerName  string
	CompanyID     uuid.UUID
	CompanyName   string
	OldRank       string
	NewRank       string
	TotalPoints   int64
	UpgradedAt    time.Time
}

// Service define las operaciones de notificación
type Service interface {
	// NotifyPointExpiring envía notificación de puntos próximos a expirar
	NotifyPointExpiring(ctx context.Context, data PointExpiringData) error

	// NotifyRedemption envía notificación de redención de puntos
	NotifyRedemption(ctx context.Context, data RedemptionData) error

	// NotifyRankUpgrade envía notificación de cambio de rango
	NotifyRankUpgrade(ctx context.Context, data RankUpgradeData) error
}
