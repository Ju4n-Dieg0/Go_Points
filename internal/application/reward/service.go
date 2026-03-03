package reward

import (
	"context"
	"errors"
	"fmt"

	subscriptionApp "github.com/Ju4n-Dieg0/Go_Points/internal/application/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/company"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/product"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/reward"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRewardNotFound         = errors.New("recompensa no encontrada")
	ErrPathNotFound           = errors.New("camino no encontrado")
	ErrInvalidCompany         = errors.New("compañía inválida")
	ErrInactiveSubscription   = errors.New("suscripción inactiva")
	ErrProductNotFound        = errors.New("producto no encontrado")
	ErrProductCompanyMismatch = errors.New("el producto no pertenece a la compañía")
	ErrRewardCompanyMismatch  = errors.New("la recompensa no pertenece a la compañía")
)

// RewardService gestiona la lógica de negocio de recompensas
type RewardService struct {
	db                  *gorm.DB
	rewardRepo          reward.RewardRepository
	pathRepo            reward.RewardPathRepository
	itemRepo            reward.RewardPathItemRepository
	companyRepo         company.Repository
	productRepo         product.Repository
	subscriptionService subscriptionApp.Service
}

// NewRewardService crea una nueva instancia
func NewRewardService(
	db *gorm.DB,
	rewardRepo reward.RewardRepository,
	pathRepo reward.RewardPathRepository,
	itemRepo reward.RewardPathItemRepository,
	companyRepo company.Repository,
	productRepo product.Repository,
	subscriptionService subscriptionApp.Service,
) *RewardService {
	return &RewardService{
		db:                  db,
		rewardRepo:          rewardRepo,
		pathRepo:            pathRepo,
		itemRepo:            itemRepo,
		companyRepo:         companyRepo,
		productRepo:         productRepo,
		subscriptionService: subscriptionService,
	}
}

// CreateReward crea una nueva recompensa
func (s *RewardService) CreateReward(ctx context.Context, companyID uuid.UUID, req *reward.CreateRewardRequest) (*reward.RewardResponse, error) {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Validar que el producto existe y pertenece a la compañía
	prod, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("error al buscar producto: %w", err)
	}

	if prod.CompanyID != companyID {
		return nil, ErrProductCompanyMismatch
	}

	// Crear recompensa
	rew := &reward.Reward{
		CompanyID:      companyID,
		ProductID:      req.ProductID,
		RequiredPoints: req.RequiredPoints,
	}

	if err := s.rewardRepo.Create(ctx, rew); err != nil {
		return nil, fmt.Errorf("error al crear recompensa: %w", err)
	}

	return reward.ToRewardResponse(rew), nil
}

// UpdateReward actualiza una recompensa
func (s *RewardService) UpdateReward(ctx context.Context, id, companyID uuid.UUID, req *reward.UpdateRewardRequest) (*reward.RewardResponse, error) {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Buscar recompensa
	rew, err := s.rewardRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRewardNotFound
		}
		return nil, fmt.Errorf("error al buscar recompensa: %w", err)
	}

	// Si se cambia el producto, validar
	if req.ProductID != nil && *req.ProductID != rew.ProductID {
		prod, err := s.productRepo.FindByID(ctx, *req.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrProductNotFound
			}
			return nil, fmt.Errorf("error al buscar producto: %w", err)
		}

		if prod.CompanyID != companyID {
			return nil, ErrProductCompanyMismatch
		}

		rew.ProductID = *req.ProductID
	}

	// Actualizar campos opcionales
	if req.RequiredPoints != nil {
		rew.RequiredPoints = *req.RequiredPoints
	}

	if err := s.rewardRepo.Update(ctx, rew); err != nil {
		return nil, fmt.Errorf("error al actualizar recompensa: %w", err)
	}

	return reward.ToRewardResponse(rew), nil
}

// DeleteReward elimina una recompensa
func (s *RewardService) DeleteReward(ctx context.Context, id, companyID uuid.UUID) error {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return err
	}

	// Verificar que existe y pertenece a la compañía
	_, err := s.rewardRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRewardNotFound
		}
		return fmt.Errorf("error al buscar recompensa: %w", err)
	}

	if err := s.rewardRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error al eliminar recompensa: %w", err)
	}

	return nil
}

// GetReward obtiene una recompensa por ID
func (s *RewardService) GetReward(ctx context.Context, id, companyID uuid.UUID) (*reward.RewardResponse, error) {
	rew, err := s.rewardRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRewardNotFound
		}
		return nil, fmt.Errorf("error al buscar recompensa: %w", err)
	}

	return reward.ToRewardResponse(rew), nil
}

// ListRewards lista recompensas de una compañía
func (s *RewardService) ListRewards(ctx context.Context, companyID uuid.UUID, page, pageSize int) (*reward.RewardListResponse, error) {
	rewards, total, err := s.rewardRepo.ListByCompany(ctx, companyID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error al listar recompensas: %w", err)
	}

	return reward.ToRewardListResponse(rewards, total, page, pageSize), nil
}

// CreateRewardPath crea un nuevo camino con sus recompensas
func (s *RewardService) CreateRewardPath(ctx context.Context, companyID uuid.UUID, req *reward.CreateRewardPathRequest) (*reward.RewardPathResponse, error) {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Validar que todas las recompensas existen y pertenecen a la compañía
	for _, rewardID := range req.RewardIDs {
		_, err := s.rewardRepo.FindByIDAndCompany(ctx, rewardID, companyID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("recompensa %s no encontrada o no pertenece a la compañía", rewardID)
			}
			return nil, fmt.Errorf("error al validar recompensa %s: %w", rewardID, err)
		}
	}

	// Crear camino y items en transacción
	var path *reward.RewardPath
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Crear camino
		path = &reward.RewardPath{
			CompanyID: companyID,
			Name:      req.Name,
		}

		if err := tx.Create(path).Error; err != nil {
			return fmt.Errorf("error al crear camino: %w", err)
		}

		// Crear items
		if len(req.RewardIDs) > 0 {
			items := make([]reward.RewardPathItem, len(req.RewardIDs))
			for i, rewardID := range req.RewardIDs {
				items[i] = reward.RewardPathItem{
					RewardPathID: path.ID,
					RewardID:     rewardID,
					Order:        i + 1,
				}
			}

			if err := tx.Create(&items).Error; err != nil {
				return fmt.Errorf("error al crear items del camino: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Obtener camino con items y recompensas
	return s.GetRewardPath(ctx, path.ID, companyID)
}

// UpdateRewardPath actualiza un camino
func (s *RewardService) UpdateRewardPath(ctx context.Context, id, companyID uuid.UUID, req *reward.UpdateRewardPathRequest) (*reward.RewardPathResponse, error) {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Buscar camino
	path, err := s.pathRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPathNotFound
		}
		return nil, fmt.Errorf("error al buscar camino: %w", err)
	}

	// Actualizar campos opcionales
	if req.Name != nil {
		path.Name = *req.Name
	}

	if err := s.pathRepo.Update(ctx, path); err != nil {
		return nil, fmt.Errorf("error al actualizar camino: %w", err)
	}

	return s.GetRewardPath(ctx, id, companyID)
}

// DeleteRewardPath elimina un camino
func (s *RewardService) DeleteRewardPath(ctx context.Context, id, companyID uuid.UUID) error {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return err
	}

	// Verificar que existe y pertenece a la compañía
	_, err := s.pathRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPathNotFound
		}
		return fmt.Errorf("error al buscar camino: %w", err)
	}

	// Eliminar en transacción
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Eliminar items
		if err := tx.Where("path_id = ?", id).Delete(&reward.RewardPathItem{}).Error; err != nil {
			return fmt.Errorf("error al eliminar items: %w", err)
		}

		// Eliminar camino
		if err := tx.Delete(&reward.RewardPath{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("error al eliminar camino: %w", err)
		}

		return nil
	})
}

// GetRewardPath obtiene un camino con sus recompensas
func (s *RewardService) GetRewardPath(ctx context.Context, id, companyID uuid.UUID) (*reward.RewardPathResponse, error) {
	// Buscar camino
	path, err := s.pathRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPathNotFound
		}
		return nil, fmt.Errorf("error al buscar camino: %w", err)
	}

	// Obtener items
	items, err := s.itemRepo.FindByPath(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener items: %w", err)
	}

	// Obtener recompensas para cada item
	itemDetails := make([]reward.PathItemDetail, 0, len(items))
	for _, item := range items {
		rew, err := s.rewardRepo.FindByID(ctx, item.RewardID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener recompensa %s: %w", item.RewardID, err)
		}

		itemDetails = append(itemDetails, reward.PathItemDetail{
			RewardID:       rew.ID,
			ProductID:      rew.ProductID,
			RequiredPoints: rew.RequiredPoints,
			Order:          item.Order,
		})
	}

	return reward.ToRewardPathResponse(path, itemDetails), nil
}

// ListRewardPaths lista caminos de una compañía
func (s *RewardService) ListRewardPaths(ctx context.Context, companyID uuid.UUID, page, pageSize int) (*reward.RewardPathListResponse, error) {
	paths, total, err := s.pathRepo.ListByCompany(ctx, companyID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error al listar caminos: %w", err)
	}

	// Obtener items para cada camino
	pathResponses := make([]reward.RewardPathResponse, 0, len(paths))
	for _, path := range paths {
		items, err := s.itemRepo.FindByPath(ctx, path.ID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener items del camino %s: %w", path.ID, err)
		}

		// Obtener recompensas para cada item
		itemDetails := make([]reward.PathItemDetail, 0, len(items))
		for _, item := range items {
			rew, err := s.rewardRepo.FindByID(ctx, item.RewardID)
			if err != nil {
				return nil, fmt.Errorf("error al obtener recompensa %s: %w", item.RewardID, err)
			}

			itemDetails = append(itemDetails, reward.PathItemDetail{
				RewardID:       rew.ID,
				ProductID:      rew.ProductID,
				RequiredPoints: rew.RequiredPoints,
				Order:          item.Order,
			})
		}

		pathResponses = append(pathResponses, *reward.ToRewardPathResponse(&path, itemDetails))
	}

	return reward.ToRewardPathListResponse(pathResponses, total, page, pageSize), nil
}

// ReorderPathItems reordena los items de un camino
func (s *RewardService) ReorderPathItems(ctx context.Context, pathID, companyID uuid.UUID, req *reward.ReorderPathItemsRequest) (*reward.RewardPathResponse, error) {
	// Validar suscripción activa
	if err := s.validateCompanySubscription(ctx, companyID); err != nil {
		return nil, err
	}

	// Verificar que el camino existe y pertenece a la compañía
	_, err := s.pathRepo.FindByIDAndCompany(ctx, pathID, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPathNotFound
		}
		return nil, fmt.Errorf("error al buscar camino: %w", err)
	}

	// Validar que todas las recompensas pertenecen a la compañía
	for _, item := range req.Items {
		_, err := s.rewardRepo.FindByIDAndCompany(ctx, item.RewardID, companyID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("recompensa %s no encontrada o no pertenece a la compañía", item.RewardID)
			}
			return nil, fmt.Errorf("error al validar recompensa %s: %w", item.RewardID, err)
		}
	}

	// Reordenar en transacción
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Eliminar items existentes
		if err := tx.Where("path_id = ?", pathID).Delete(&reward.RewardPathItem{}).Error; err != nil {
			return fmt.Errorf("error al eliminar items existentes: %w", err)
		}

		// Crear nuevos items
		if len(req.Items) > 0 {
			items := make([]reward.RewardPathItem, len(req.Items))
			for i, item := range req.Items {
				items[i] = reward.RewardPathItem{
					RewardPathID: pathID,
					RewardID: item.RewardID,
					Order:    item.Order,
				}
			}

			if err := tx.Create(&items).Error; err != nil {
				return fmt.Errorf("error al crear nuevos items: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetRewardPath(ctx, pathID, companyID)
}

// validateCompanySubscription valida que la compañía existe y tiene suscripción activa
func (s *RewardService) validateCompanySubscription(ctx context.Context, companyID uuid.UUID) error {
	// Verificar que la compañía existe
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidCompany
		}
		return fmt.Errorf("error al buscar compañía: %w", err)
	}

	// Verificar suscripción activa
	err = s.subscriptionService.ValidateActiveSubscription(ctx, companyID)
	if err != nil {
		return ErrInactiveSubscription
	}

	return nil
}
