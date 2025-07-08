package services

import (
	"context"
	"example.com/smart-devices/internal/models"
	"fmt"
	"go.uber.org/zap"
)

// DeviceRepository is the minimal interface DeviceService needs.
// Both *repository.DeviceRepository and *MockDeviceRepository satisfy this.
type DeviceRepository interface {
	GetDevice(ctx context.Context, id string) (*models.Device, error)
	GetDevices(ctx context.Context) ([]models.Device, error)
	CreateDevice(ctx context.Context, device models.Device) (models.Device, error)
	UpdateDevice(ctx context.Context, id string, device models.Device) (*models.Device, error)
	DeleteDevice(ctx context.Context, id string) error
	UpdateDeviceHomeID(ctx context.Context, id, homeID string) error
}

type DeviceService struct {
	repo   DeviceRepository
	logger *zap.Logger
}

// NewDeviceService accepts any DeviceRepository (mock or real).
func NewDeviceService(repo DeviceRepository, logger *zap.Logger) *DeviceService {
	return &DeviceService{
		repo:   repo,
		logger: logger,
	}
}

func (s *DeviceService) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	s.logger.Debug("fetching device",
		zap.String("device_id", id),
		zap.String("layer", "service"),
	)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	device, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		s.logger.Warn("device retrieval failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return device, nil
}

func (s *DeviceService) GetDevices(ctx context.Context) ([]models.Device, error) {
	s.logger.Debug("fetching device",
		zap.String("layer", "service"),
	)

	device, err := s.repo.GetDevices(ctx)
	if err != nil {
		s.logger.Warn("device retrieval failed",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return device, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	s.logger.Debug("deleting device",
		zap.String("device_id", id),
		zap.String("layer", "service"),
	)

	if id == "" {
		return fmt.Errorf("id is required")
	}

	err := s.repo.DeleteDevice(ctx, id)
	if err != nil {
		s.logger.Warn("device deletion failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get device: %w", err)
	}

	return nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id string, device models.Device) (*models.Device, error) {
	s.logger.Debug("updating device",
		zap.String("device_id", id),
		zap.String("layer", "service"),
	)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	updatedDevice, err := s.repo.UpdateDevice(ctx, id, device)
	if err != nil {
		s.logger.Warn("device update failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return updatedDevice, nil
}

func (s *DeviceService) CreateDevice(ctx context.Context, device models.Device) (models.Device, error) {
	s.logger.Debug("creating device",
		zap.String("device_id", device.ID),
		zap.String("layer", "service"),
	)

	device, err := s.repo.CreateDevice(ctx, device)
	if err != nil {
		s.logger.Warn("device creation failed",
			zap.Error(err),
		)
		return device, fmt.Errorf("failed to update device: %w", err)
	}

	return device, nil
}

func (s *DeviceService) UpdateDeviceHomeID(ctx context.Context, id string, homeID string) error {
	s.logger.Debug("updating device home id",
		zap.String("device_id", id),
		zap.String("layer", "service"),
	)

	if id == "" {
		return fmt.Errorf("id is required")
	}

	err := s.repo.UpdateDeviceHomeID(ctx, id, homeID)
	if err != nil {
		s.logger.Warn("device update failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update device: %w", err)
	}

	return nil
}
