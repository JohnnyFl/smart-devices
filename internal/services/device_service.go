package services

import (
	"context"
	"example.com/smart-devices/internal/errors"
	"example.com/smart-devices/internal/models"
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
		return nil, errors.ErrDomainInvalidDeviceID.
			WithOperation("GetDevice").
			WithLayer("service").
			WithContext("reason", "device ID is empty")
	}

	device, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		// Check if it's already a domain error and preserve it
		if domainErr, ok := err.(*errors.DomainError); ok {
			s.logger.Warn("device retrieval failed",
				zap.String("device_id", id),
				zap.String("error_type", string(domainErr.Type)),
				zap.Error(err),
			)
			return nil, domainErr.WithLayer("service")
		}

		// Wrap unknown errors
		s.logger.Warn("device retrieval failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeInternal, "failed to retrieve device", err).
			WithOperation("GetDevice").
			WithLayer("service").
			WithContext("device_id", id)
	}

	return device, nil
}

func (s *DeviceService) GetDevices(ctx context.Context) ([]models.Device, error) {
	s.logger.Debug("fetching devices",
		zap.String("layer", "service"),
	)

	devices, err := s.repo.GetDevices(ctx)
	if err != nil {
		// Check if it's already a domain error and preserve it
		if domainErr, ok := err.(*errors.DomainError); ok {
			s.logger.Warn("devices retrieval failed",
				zap.String("error_type", string(domainErr.Type)),
				zap.Error(err),
			)
			return nil, domainErr.WithLayer("service")
		}

		// Wrap unknown errors
		s.logger.Warn("devices retrieval failed",
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeInternal, "failed to retrieve devices", err).
			WithOperation("GetDevices").
			WithLayer("service")
	}

	return devices, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	s.logger.Debug("deleting device",
		zap.String("device_id", id),
		zap.String("layer", "service"),
	)

	if id == "" {
		return errors.ErrDomainInvalidDeviceID.
			WithOperation("DeleteDevice").
			WithLayer("service").
			WithContext("reason", "device ID is empty")
	}

	err := s.repo.DeleteDevice(ctx, id)
	if err != nil {
		// Check if it's already a domain error and preserve it
		if domainErr, ok := err.(*errors.DomainError); ok {
			s.logger.Warn("device deletion failed",
				zap.String("device_id", id),
				zap.String("error_type", string(domainErr.Type)),
				zap.Error(err),
			)
			return domainErr.WithLayer("service")
		}

		// Wrap unknown errors
		s.logger.Warn("device deletion failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrorTypeInternal, "failed to delete device", err).
			WithOperation("DeleteDevice").
			WithLayer("service").
			WithContext("device_id", id)
	}

	return nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id string, device models.Device) (*models.Device, error) {
	s.logger.Debug("updating device",
		zap.String("device_id", id),
		zap.String("layer", "service"),
	)

	if id == "" {
		return nil, errors.ErrDomainInvalidDeviceID.
			WithOperation("UpdateDevice").
			WithLayer("service").
			WithContext("reason", "device ID is empty")
	}

	updatedDevice, err := s.repo.UpdateDevice(ctx, id, device)
	if err != nil {
		// Check if it's already a domain error and preserve it
		if domainErr, ok := err.(*errors.DomainError); ok {
			s.logger.Warn("device update failed",
				zap.String("device_id", id),
				zap.String("error_type", string(domainErr.Type)),
				zap.Error(err),
			)
			return nil, domainErr.WithLayer("service")
		}

		// Wrap unknown errors
		s.logger.Warn("device update failed",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeInternal, "failed to update device", err).
			WithOperation("UpdateDevice").
			WithLayer("service").
			WithContext("device_id", id)
	}

	return updatedDevice, nil
}

func (s *DeviceService) CreateDevice(ctx context.Context, device models.Device) (models.Device, error) {
	s.logger.Debug("creating device",
		zap.String("device_mac", device.MAC),
		zap.String("device_name", device.Name),
		zap.String("layer", "service"),
	)

	createdDevice, err := s.repo.CreateDevice(ctx, device)
	if err != nil {
		// Check if it's already a domain error and preserve it
		if domainErr, ok := err.(*errors.DomainError); ok {
			s.logger.Warn("device creation failed",
				zap.String("device_mac", device.MAC),
				zap.String("error_type", string(domainErr.Type)),
				zap.Error(err),
			)
			return device, domainErr.WithLayer("service")
		}

		// Wrap unknown errors
		s.logger.Warn("device creation failed",
			zap.String("device_mac", device.MAC),
			zap.Error(err),
		)
		return device, errors.WrapError(errors.ErrorTypeInternal, "failed to create device", err).
			WithOperation("CreateDevice").
			WithLayer("service").
			WithContext("device_mac", device.MAC)
	}

	return createdDevice, nil
}

func (s *DeviceService) UpdateDeviceHomeID(ctx context.Context, id string, homeID string) error {
	s.logger.Debug("updating device home id",
		zap.String("device_id", id),
		zap.String("home_id", homeID),
		zap.String("layer", "service"),
	)

	if id == "" {
		return errors.ErrDomainInvalidDeviceID.
			WithOperation("UpdateDeviceHomeID").
			WithLayer("service").
			WithContext("reason", "device ID is empty")
	}

	if homeID == "" {
		return errors.ErrDomainMissingHomeID.
			WithOperation("UpdateDeviceHomeID").
			WithLayer("service").
			WithContext("device_id", id)
	}

	err := s.repo.UpdateDeviceHomeID(ctx, id, homeID)
	if err != nil {
		// Check if it's already a domain error and preserve it
		if domainErr, ok := err.(*errors.DomainError); ok {
			s.logger.Warn("device home ID update failed",
				zap.String("device_id", id),
				zap.String("home_id", homeID),
				zap.String("error_type", string(domainErr.Type)),
				zap.Error(err),
			)
			return domainErr.WithLayer("service")
		}

		// Wrap unknown errors
		s.logger.Warn("device home ID update failed",
			zap.String("device_id", id),
			zap.String("home_id", homeID),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrorTypeInternal, "failed to update device home ID", err).
			WithOperation("UpdateDeviceHomeID").
			WithLayer("service").
			WithContext("device_id", id).
			WithContext("home_id", homeID)
	}

	return nil
}
