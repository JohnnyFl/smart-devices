package services

import (
	"context"
	"encoding/json"
	"example.com/smart-devices/internal/models"
	"go.uber.org/zap"
)

type SQSService struct {
	deviceService *DeviceService
	logger        *zap.Logger
}

func NewSQSService(deviceService *DeviceService, logger *zap.Logger) *SQSService {
	return &SQSService{
		deviceService: deviceService,
		logger:        logger,
	}
}

func (s *SQSService) ProcessMessage(ctx context.Context, msg string) error {
	var message models.SQSMessage

	if err := json.Unmarshal([]byte(msg), &message); err != nil {
		s.logger.Error("failed to unmarshal message", zap.Error(err))
		return err
	}

	s.logger.Info("processing device-home association", zap.String("device-id", message.DeviceID), zap.String("home-id", message.HomeID))

	if err := s.deviceService.UpdateDeviceHomeID(ctx, message.DeviceID, message.HomeID); err != nil {
		s.logger.Error("failed to update device-home association", zap.Error(err), zap.String("device-id", message.DeviceID), zap.String("home-id", message.HomeID))
		return err
	}
	s.logger.Info("device-home association updated", zap.String("device-id", message.DeviceID), zap.String("home-id", message.HomeID))
	return nil

}
