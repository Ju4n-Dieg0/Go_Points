package point

import (
	"context"
	"fmt"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	domainPoint "github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service define la interfaz para operaciones de puntos
type Service interface {
	// EarnPoints registra puntos ganados
	EarnPoints(ctx context.Context, req *domainPoint.EarnPointsRequest) (*domainPoint.PointBalanceResponse, error)

	// RedeemPoints redime puntos usando FIFO
	RedeemPoints(ctx context.Context, req *domainPoint.RedeemPointsRequest) (*domainPoint.PointBalanceResponse, error)

	// GetBalance obtiene el balance de puntos
	GetBalance(ctx context.Context, consumerID, companyID uuid.UUID) (*domainPoint.PointBalanceResponse, error)

	// GetTransactions obtiene el historial de transacciones
	GetTransactions(ctx context.Context, consumerID, companyID uuid.UUID, page, pageSize int) (*domainPoint.TransactionListResponse, error)

	// ConfigureRank configura los rangos de una empresa
	ConfigureRank(ctx context.Context, companyID uuid.UUID, req *domainPoint.ConfigureRankRequest) (*domainPoint.RankConfigResponse, error)

	// GetRankConfig obtiene la configuración de rangos
	GetRankConfig(ctx context.Context, companyID uuid.UUID) (*domainPoint.RankConfigResponse, error)

	// ProcessExpiredPoints procesa puntos expirados (cronjob)
	ProcessExpiredPoints(ctx context.Context) error

	// ProcessInactivityPenalties procesa penalizaciones por inactividad (cronjob)
	ProcessInactivityPenalties(ctx context.Context) error

	// NotifyExpiringSoon notifica puntos que expirarán pronto (cronjob)
	NotifyExpiringSoon(ctx context.Context) error
}

// service implementa Service
type service struct {
	db                  *gorm.DB
	balanceRepo         domainPoint.BalanceRepository
	transactionRepo     domainPoint.TransactionRepository
	rankConfigRepo      domainPoint.RankConfigRepository
	notificationService domainPoint.NotificationService
	config              config.PointsConfig
}

// NewService crea una nueva instancia de Service
func NewService(
	db *gorm.DB,
	balanceRepo domainPoint.BalanceRepository,
	transactionRepo domainPoint.TransactionRepository,
	rankConfigRepo domainPoint.RankConfigRepository,
	notificationService domainPoint.NotificationService,
	cfg config.PointsConfig,
) Service {
	return &service{
		db:                  db,
		balanceRepo:         balanceRepo,
		transactionRepo:     transactionRepo,
		rankConfigRepo:      rankConfigRepo,
		notificationService: notificationService,
		config:              cfg,
	}
}

// EarnPoints registra puntos ganados
func (s *service) EarnPoints(ctx context.Context, req *domainPoint.EarnPointsRequest) (*domainPoint.PointBalanceResponse, error) {
	var response *domainPoint.PointBalanceResponse

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Buscar o crear balance
		balance, err := s.balanceRepo.FindByConsumerAndCompanyForUpdate(ctx, tx, req.ConsumerID, req.CompanyID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Crear nuevo balance
				balance = &domainPoint.ConsumerCompanyPoints{
					ID:                    uuid.New(),
					ConsumerID:            req.ConsumerID,
					CompanyID:             req.CompanyID,
					TotalHistoricalPoints: 0,
					TotalAvailablePoints:  0,
				}
				if err := s.balanceRepo.Create(ctx, tx, balance); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Actualizar balance
		balance.TotalHistoricalPoints += req.Points
		balance.TotalAvailablePoints += req.Points

		if err := s.balanceRepo.Update(ctx, tx, balance); err != nil {
			return err
		}

		// Crear transacción con fecha de expiración
		expirationDate := time.Now().AddDate(0, s.config.ExpirationMonths, 0)
		transaction := &domainPoint.PointTransaction{
			ID:              uuid.New(),
			ConsumerID:      req.ConsumerID,
			CompanyID:       req.CompanyID,
			Points:          req.Points,
			RemainingPoints: req.Points,
			Type:            domainPoint.TransactionTypeEarn,
			ExpirationDate:  &expirationDate,
		}

		if err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
			return err
		}

		// Obtener configuración de rangos
		rankConfig, _ := s.rankConfigRepo.FindByCompany(ctx, req.CompanyID)
		rank := balance.GetRank(rankConfig)

		response = &domainPoint.PointBalanceResponse{
			ConsumerID:            balance.ConsumerID,
			CompanyID:             balance.CompanyID,
			TotalHistoricalPoints: balance.TotalHistoricalPoints,
			TotalAvailablePoints:  balance.TotalAvailablePoints,
			Rank:                  rank,
			LastRedemptionDate:    balance.LastRedemptionDate,
			CreatedAt:             balance.CreatedAt,
			UpdatedAt:             balance.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		logger.Error("Failed to earn points", "error", err)
		return nil, errors.ErrDatabase.WithError(err)
	}

	logger.Info("Points earned successfully",
		"consumer_id", req.ConsumerID,
		"company_id", req.CompanyID,
		"points", req.Points,
	)

	return response, nil
}

// RedeemPoints redime puntos usando FIFO real
func (s *service) RedeemPoints(ctx context.Context, req *domainPoint.RedeemPointsRequest) (*domainPoint.PointBalanceResponse, error) {
	var response *domainPoint.PointBalanceResponse

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Buscar balance con lock
		balance, err := s.balanceRepo.FindByConsumerAndCompanyForUpdate(ctx, tx, req.ConsumerID, req.CompanyID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.ErrNotFound.WithError(fmt.Errorf("balance not found"))
			}
			return err
		}

		// Validar puntos disponibles
		if balance.TotalAvailablePoints < req.Points {
			return errors.ErrValidation.WithError(
				fmt.Errorf("insufficient points: available=%d, requested=%d",
					balance.TotalAvailablePoints, req.Points),
			)
		}

		// Obtener transacciones disponibles para redimir en orden FIFO
		// CRÍTICO: SELECT FOR UPDATE para evitar race conditions
		availableTransactions, err := s.transactionRepo.FindAvailableForRedemption(ctx, tx, req.ConsumerID, req.CompanyID)
		if err != nil {
			return err
		}

		if len(availableTransactions) == 0 {
			return errors.ErrValidation.WithError(fmt.Errorf("no available transactions to redeem"))
		}

		// Aplicar FIFO: redimir puntos de las transacciones más antiguas primero
		pointsToRedeem := req.Points
		updatedTransactions := []domainPoint.PointTransaction{}

		for i := range availableTransactions {
			if pointsToRedeem <= 0 {
				break
			}

			tx := &availableTransactions[i]

			if tx.RemainingPoints <= 0 {
				continue
			}

			// Determinar cuántos puntos redimir de esta transacción
			pointsFromThisTx := tx.RemainingPoints
			if pointsFromThisTx > pointsToRedeem {
				pointsFromThisTx = pointsToRedeem
			}

			// Actualizar remaining points
			tx.RemainingPoints -= pointsFromThisTx
			pointsToRedeem -= pointsFromThisTx

			updatedTransactions = append(updatedTransactions, *tx)
		}

		// Guardar transacciones actualizadas
		if err := s.transactionRepo.UpdateBatch(ctx, tx, updatedTransactions); err != nil {
			return err
		}

		// Crear transacción de redención
		redeemTransaction := &domainPoint.PointTransaction{
			ID:              uuid.New(),
			ConsumerID:      req.ConsumerID,
			CompanyID:       req.CompanyID,
			Points:          -req.Points, // Negativo para redención
			RemainingPoints: 0,
			Type:            domainPoint.TransactionTypeRedeem,
		}

		if err := s.transactionRepo.Create(ctx, tx, redeemTransaction); err != nil {
			return err
		}

		// Actualizar balance
		balance.TotalAvailablePoints -= req.Points
		now := time.Now()
		balance.LastRedemptionDate = &now

		if err := s.balanceRepo.Update(ctx, tx, balance); err != nil {
			return err
		}

		// Obtener configuración de rangos
		rankConfig, _ := s.rankConfigRepo.FindByCompany(ctx, req.CompanyID)
		rank := balance.GetRank(rankConfig)

		response = &domainPoint.PointBalanceResponse{
			ConsumerID:            balance.ConsumerID,
			CompanyID:             balance.CompanyID,
			TotalHistoricalPoints: balance.TotalHistoricalPoints,
			TotalAvailablePoints:  balance.TotalAvailablePoints,
			Rank:                  rank,
			LastRedemptionDate:    balance.LastRedemptionDate,
			CreatedAt:             balance.CreatedAt,
			UpdatedAt:             balance.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		logger.Error("Failed to redeem points", "error", err)
		return nil, err
	}

	logger.Info("Points redeemed successfully",
		"consumer_id", req.ConsumerID,
		"company_id", req.CompanyID,
		"points", req.Points,
	)

	return response, nil
}

// GetBalance obtiene el balance de puntos
func (s *service) GetBalance(ctx context.Context, consumerID, companyID uuid.UUID) (*domainPoint.PointBalanceResponse, error) {
	balance, err := s.balanceRepo.FindByConsumerAndCompany(ctx, consumerID, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithError(fmt.Errorf("balance not found"))
		}
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Obtener configuración de rangos
	rankConfig, _ := s.rankConfigRepo.FindByCompany(ctx, companyID)
	rank := balance.GetRank(rankConfig)

	response := domainPoint.ToPointBalanceResponse(balance, rank)
	return &response, nil
}

// GetTransactions obtiene el historial de transacciones
func (s *service) GetTransactions(ctx context.Context, consumerID, companyID uuid.UUID, page, pageSize int) (*domainPoint.TransactionListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	transactions, total, err := s.transactionRepo.FindByConsumerAndCompany(ctx, consumerID, companyID, page, pageSize)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainPoint.ToTransactionListResponse(transactions, total, page, pageSize)
	return &response, nil
}

// ConfigureRank configura los rangos de una empresa
func (s *service) ConfigureRank(ctx context.Context, companyID uuid.UUID, req *domainPoint.ConfigureRankRequest) (*domainPoint.RankConfigResponse, error) {
	config := &domainPoint.CompanyRankConfig{
		ID:              uuid.New(),
		CompanyID:       companyID,
		SilverMinPoints: req.SilverMinPoints,
		GoldMinPoints:   req.GoldMinPoints,
	}

	// Validar
	if err := config.Validate(); err != nil {
		return nil, errors.ErrValidation.WithError(err)
	}

	// Upsert
	if err := s.rankConfigRepo.Upsert(ctx, config); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainPoint.ToRankConfigResponse(config)
	return &response, nil
}

// GetRankConfig obtiene la configuración de rangos
func (s *service) GetRankConfig(ctx context.Context, companyID uuid.UUID) (*domainPoint.RankConfigResponse, error) {
	config, err := s.rankConfigRepo.FindByCompany(ctx, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithError(fmt.Errorf("rank configuration not found"))
		}
		return nil, errors.ErrDatabase.WithError(err)
	}

	response := domainPoint.ToRankConfigResponse(config)
	return &response, nil
}

// ProcessExpiredPoints procesa puntos expirados
func (s *service) ProcessExpiredPoints(ctx context.Context) error {
	expiredTransactions, err := s.transactionRepo.FindExpired(ctx)
	if err != nil {
		logger.Error("Failed to find expired transactions", "error", err)
		return err
	}

	logger.Info("Processing expired points", "count", len(expiredTransactions))

	for _, tx := range expiredTransactions {
		if err := s.expireTransaction(ctx, &tx); err != nil {
			logger.Error("Failed to expire transaction",
				"transaction_id", tx.ID,
				"error", err,
			)
			continue
		}
	}

	return nil
}

// expireTransaction procesa una transacción expirada
func (s *service) expireTransaction(ctx context.Context, tx *domainPoint.PointTransaction) error {
	return s.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		// Buscar balance con lock
		balance, err := s.balanceRepo.FindByConsumerAndCompanyForUpdate(ctx, dbTx, tx.ConsumerID, tx.CompanyID)
		if err != nil {
			return err
		}

		pointsToExpire := tx.RemainingPoints

		// Actualizar transacción original
		tx.RemainingPoints = 0
		if err := s.transactionRepo.UpdateBatch(ctx, dbTx, []domainPoint.PointTransaction{*tx}); err != nil {
			return err
		}

		// Crear transacción de expiración
		expireTransaction := &domainPoint.PointTransaction{
			ID:              uuid.New(),
			ConsumerID:      tx.ConsumerID,
			CompanyID:       tx.CompanyID,
			Points:          -pointsToExpire,
			RemainingPoints: 0,
			Type:            domainPoint.TransactionTypeExpire,
		}

		if err := s.transactionRepo.Create(ctx, dbTx, expireTransaction); err != nil {
			return err
		}

		// Actualizar balance
		balance.TotalAvailablePoints -= pointsToExpire
		if err := s.balanceRepo.Update(ctx, dbTx, balance); err != nil {
			return err
		}

		// Notificar
		_ = s.notificationService.NotifyPointsExpired(ctx, tx.ConsumerID, tx.CompanyID, pointsToExpire)

		logger.Info("Points expired",
			"consumer_id", tx.ConsumerID,
			"company_id", tx.CompanyID,
			"points", pointsToExpire,
		)

		return nil
	})
}

// ProcessInactivityPenalties procesa penalizaciones por inactividad
func (s *service) ProcessInactivityPenalties(ctx context.Context) error {
	inactiveBalances, err := s.balanceRepo.FindInactiveConsumers(ctx, s.config.InactivityMonths)
	if err != nil {
		logger.Error("Failed to find inactive consumers", "error", err)
		return err
	}

	logger.Info("Processing inactivity penalties", "count", len(inactiveBalances))

	for _, balance := range inactiveBalances {
		if err := s.applyInactivityPenalty(ctx, &balance); err != nil {
			logger.Error("Failed to apply inactivity penalty",
				"consumer_id", balance.ConsumerID,
				"company_id", balance.CompanyID,
				"error", err,
			)
			continue
		}
	}

	return nil
}

// applyInactivityPenalty aplica penalización por inactividad
func (s *service) applyInactivityPenalty(ctx context.Context, balance *domainPoint.ConsumerCompanyPoints) error {
	penalty := s.config.InactivityPenalty

	// No penalizar si no tiene puntos suficientes
	if balance.TotalAvailablePoints < penalty {
		penalty = balance.TotalAvailablePoints
	}

	if penalty <= 0 {
		return nil
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Actualizar balance
		balance.TotalAvailablePoints -= penalty
		if err := s.balanceRepo.Update(ctx, tx, balance); err != nil {
			return err
		}

		// Crear transacción de penalización
		penaltyTransaction := &domainPoint.PointTransaction{
			ID:              uuid.New(),
			ConsumerID:      balance.ConsumerID,
			CompanyID:       balance.CompanyID,
			Points:          -penalty,
			RemainingPoints: 0,
			Type:            domainPoint.TransactionTypePenalty,
		}

		if err := s.transactionRepo.Create(ctx, tx, penaltyTransaction); err != nil {
			return err
		}

		// Notificar
		_ = s.notificationService.NotifyInactivityPenalty(ctx, balance.ConsumerID, balance.CompanyID, penalty)

		logger.Info("Inactivity penalty applied",
			"consumer_id", balance.ConsumerID,
			"company_id", balance.CompanyID,
			"penalty", penalty,
		)

		return nil
	})
}

// NotifyExpiringSoon notifica puntos que expirarán pronto
func (s *service) NotifyExpiringSoon(ctx context.Context) error {
	transactions, err := s.transactionRepo.FindExpiringSoon(ctx, s.config.NotificationBeforeDays)
	if err != nil {
		logger.Error("Failed to find expiring transactions", "error", err)
		return err
	}

	logger.Info("Notifying expiring points", "count", len(transactions))

	for _, tx := range transactions {
		if tx.ExpirationDate == nil {
			continue
		}

		if err := s.notificationService.NotifyPointsExpiring(
			ctx,
			tx.ConsumerID,
			tx.CompanyID,
			tx.RemainingPoints,
			*tx.ExpirationDate,
		); err != nil {
			logger.Error("Failed to send notification",
				"consumer_id", tx.ConsumerID,
				"error", err,
			)
		}
	}

	return nil
}
