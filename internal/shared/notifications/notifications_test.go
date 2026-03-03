package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockNotificationService es un mock para testing
type MockNotificationService struct {
	PointExpiringCalled bool
	RedemptionCalled    bool
	RankUpgradeCalled   bool
	
	LastPointExpiringData PointExpiringData
	LastRedemptionData    RedemptionData
	LastRankUpgradeData   RankUpgradeData
	
	ShouldReturnError bool
	ErrorToReturn     error
}

func (m *MockNotificationService) NotifyPointExpiring(ctx context.Context, data PointExpiringData) error {
	m.PointExpiringCalled = true
	m.LastPointExpiringData = data
	if m.ShouldReturnError {
		return m.ErrorToReturn
	}
	return nil
}

func (m *MockNotificationService) NotifyRedemption(ctx context.Context, data RedemptionData) error {
	m.RedemptionCalled = true
	m.LastRedemptionData = data
	if m.ShouldReturnError {
		return m.ErrorToReturn
	}
	return nil
}

func (m *MockNotificationService) NotifyRankUpgrade(ctx context.Context, data RankUpgradeData) error {
	m.RankUpgradeCalled = true
	m.LastRankUpgradeData = data
	if m.ShouldReturnError {
		return m.ErrorToReturn
	}
	return nil
}

func TestPointExpiringNotification(t *testing.T) {
	ctx := context.Background()
	mock := &MockNotificationService{}
	
	consumerID := uuid.New()
	companyID := uuid.New()
	expirationDate := time.Now().AddDate(0, 0, 7)
	
	data := PointExpiringData{
		ConsumerID:      consumerID,
		ConsumerEmail:   "test@example.com",
		ConsumerName:    "Test User",
		CompanyID:       companyID,
		CompanyName:     "Test Company",
		Points:          500,
		ExpirationDate:  expirationDate,
		DaysUntilExpiry: 7,
	}
	
	err := mock.NotifyPointExpiring(ctx, data)
	
	require.NoError(t, err)
	assert.True(t, mock.PointExpiringCalled)
	assert.Equal(t, data.ConsumerID, mock.LastPointExpiringData.ConsumerID)
	assert.Equal(t, data.Points, mock.LastPointExpiringData.Points)
}

func TestRedemptionNotification(t *testing.T) {
	ctx := context.Background()
	mock := &MockNotificationService{}
	
	data := RedemptionData{
		ConsumerID:       uuid.New(),
		ConsumerEmail:    "test@example.com",
		ConsumerName:     "Test User",
		CompanyID:        uuid.New(),
		CompanyName:      "Test Company",
		PointsRedeemed:   100,
		RemainingBalance: 400,
		TransactionID:    uuid.New(),
		RedeemedAt:       time.Now(),
	}
	
	err := mock.NotifyRedemption(ctx, data)
	
	require.NoError(t, err)
	assert.True(t, mock.RedemptionCalled)
	assert.Equal(t, int64(100), mock.LastRedemptionData.PointsRedeemed)
	assert.Equal(t, int64(400), mock.LastRedemptionData.RemainingBalance)
}

func TestRankUpgradeNotification(t *testing.T) {
	ctx := context.Background()
	mock := &MockNotificationService{}
	
	data := RankUpgradeData{
		ConsumerID:    uuid.New(),
		ConsumerEmail: "test@example.com",
		ConsumerName:  "Test User",
		CompanyID:     uuid.New(),
		CompanyName:   "Test Company",
		OldRank:       "SILVER",
		NewRank:       "GOLD",
		TotalPoints:   10000,
		UpgradedAt:    time.Now(),
	}
	
	err := mock.NotifyRankUpgrade(ctx, data)
	
	require.NoError(t, err)
	assert.True(t, mock.RankUpgradeCalled)
	assert.Equal(t, "SILVER", mock.LastRankUpgradeData.OldRank)
	assert.Equal(t, "GOLD", mock.LastRankUpgradeData.NewRank)
}

func TestCompositeNotificationService(t *testing.T) {
	ctx := context.Background()
	
	mock1 := &MockNotificationService{}
	mock2 := &MockNotificationService{}
	
	composite := NewCompositeNotificationService(mock1, mock2)
	
	data := PointExpiringData{
		ConsumerID:      uuid.New(),
		ConsumerEmail:   "test@example.com",
		ConsumerName:    "Test User",
		CompanyID:       uuid.New(),
		CompanyName:     "Test Company",
		Points:          500,
		ExpirationDate:  time.Now().AddDate(0, 0, 7),
		DaysUntilExpiry: 7,
	}
	
	err := composite.NotifyPointExpiring(ctx, data)
	
	require.NoError(t, err)
	assert.True(t, mock1.PointExpiringCalled, "First service should be called")
	assert.True(t, mock2.PointExpiringCalled, "Second service should be called")
}
