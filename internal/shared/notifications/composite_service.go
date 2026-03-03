package notifications

import (
	"context"
	"sync"
)

// CompositeNotificationService envía notificaciones a múltiples servicios
type CompositeNotificationService struct {
	services []Service
}

// NewCompositeNotificationService crea un servicio que envía a múltiples destinos
func NewCompositeNotificationService(services ...Service) *CompositeNotificationService {
	return &CompositeNotificationService{
		services: services,
	}
}

// NotifyPointExpiring envía la notificación a todos los servicios configurados
func (s *CompositeNotificationService) NotifyPointExpiring(ctx context.Context, data PointExpiringData) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.services))

	for _, service := range s.services {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			if err := svc.NotifyPointExpiring(ctx, data); err != nil {
				errChan <- err
			}
		}(service)
	}

	wg.Wait()
	close(errChan)

	// Retornar el primer error encontrado (si existe)
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// NotifyRedemption envía la notificación a todos los servicios configurados
func (s *CompositeNotificationService) NotifyRedemption(ctx context.Context, data RedemptionData) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.services))

	for _, service := range s.services {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			if err := svc.NotifyRedemption(ctx, data); err != nil {
				errChan <- err
			}
		}(service)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// NotifyRankUpgrade envía la notificación a todos los servicios configurados
func (s *CompositeNotificationService) NotifyRankUpgrade(ctx context.Context, data RankUpgradeData) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.services))

	for _, service := range s.services {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			if err := svc.NotifyRankUpgrade(ctx, data); err != nil {
				errChan <- err
			}
		}(service)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
