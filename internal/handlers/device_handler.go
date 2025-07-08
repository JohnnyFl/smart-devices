package handlers

import (
	"context"
	"example.com/smart-devices/internal/errors"
	"example.com/smart-devices/internal/models"
	"example.com/smart-devices/internal/services"
	"example.com/smart-devices/internal/validation"
	"example.com/smart-devices/utils"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

type DeviceHandler struct {
	svc    *services.DeviceService
	logger *zap.Logger
}

func NewDeviceHandler(svc *services.DeviceService, logger *zap.Logger) *DeviceHandler {
	return &DeviceHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *DeviceHandler) GetDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deviceID, ok := request.PathParameters["id"]
	if !ok || deviceID == "" {
		return errors.ErrMissingDeviceID.ToResponse(), nil
	}

	// Validate device ID format
	if err := validation.ValidateDeviceID(deviceID); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	h.logger.Debug("fetching device",
		zap.String("device_id", deviceID),
		zap.String("layer", "handler"),
	)

	device, err := h.svc.GetDevice(ctx, deviceID)
	if err != nil {
		// Check if it's a domain error and convert appropriately
		if domainErr, ok := err.(*errors.DomainError); ok {
			h.logger.Warn("device retrieval failed",
				zap.String("device_id", deviceID),
				zap.String("error_type", string(domainErr.Type)),
				zap.String("operation", domainErr.Operation),
				zap.Error(err),
			)
			return domainErr.ToAPIError().ToResponse(), nil
		}

		// Fallback for unknown errors
		h.logger.Error("unexpected error during device retrieval",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return errors.ErrInternalServer.ToResponse(), nil
	}

	return utils.JSONSuccessResponse(200, device), nil
}

func (h *DeviceHandler) GetDevices(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	devices, err := h.svc.GetDevices(ctx)
	if err != nil {
		// Check if it's a domain error and convert appropriately
		if domainErr, ok := err.(*errors.DomainError); ok {
			h.logger.Warn("devices retrieval failed",
				zap.String("error_type", string(domainErr.Type)),
				zap.String("operation", domainErr.Operation),
				zap.Error(err),
			)
			return domainErr.ToAPIError().ToResponse(), nil
		}

		// Fallback for unknown errors
		h.logger.Error("unexpected error during devices retrieval",
			zap.Error(err),
		)
		return errors.ErrInternalServer.ToResponse(), nil
	}

	return utils.JSONSuccessResponse(200, devices), nil
}

func (h *DeviceHandler) DeleteDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deviceID, ok := request.PathParameters["id"]
	if !ok || deviceID == "" {
		return errors.ErrMissingDeviceID.ToResponse(), nil
	}

	// Validate device ID format
	if err := validation.ValidateDeviceID(deviceID); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	h.logger.Debug("deleting device",
		zap.String("device_id", deviceID),
		zap.String("layer", "handler"),
	)

	err := h.svc.DeleteDevice(ctx, deviceID)
	if err != nil {
		// Check if it's a domain error and convert appropriately
		if domainErr, ok := err.(*errors.DomainError); ok {
			h.logger.Warn("device deletion failed",
				zap.String("device_id", deviceID),
				zap.String("error_type", string(domainErr.Type)),
				zap.String("operation", domainErr.Operation),
				zap.Error(err),
			)
			return domainErr.ToAPIError().ToResponse(), nil
		}

		// Fallback for unknown errors
		h.logger.Error("unexpected error during device deletion",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return errors.ErrDeviceDeletionFailed.ToResponse(), nil
	}

	return utils.JSONSuccessResponse(200, map[string]string{"message": "Device deleted successfully"}), nil
}

func (h *DeviceHandler) UpdateDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deviceID, ok := request.PathParameters["id"]
	if !ok || deviceID == "" {
		return errors.ErrMissingDeviceID.ToResponse(), nil
	}

	// Validate device ID format
	if err := validation.ValidateDeviceID(deviceID); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	// Validate and parse request body
	var updateReq models.UpdateDeviceRequest
	if err := validation.ValidateJSON(request.Body, &updateReq); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	// Validate request data
	if err := validation.ValidateUpdateDeviceRequest(updateReq); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	// Convert to Device model for service layer
	device := models.Device{}
	if updateReq.Name != nil {
		device.Name = *updateReq.Name
	}
	if updateReq.Type != nil {
		device.Type = *updateReq.Type
	}
	if updateReq.HomeID != nil {
		device.HomeID = *updateReq.HomeID
	}

	h.logger.Debug("updating device",
		zap.String("device_id", deviceID),
		zap.String("layer", "handler"),
	)

	updatedDevice, err := h.svc.UpdateDevice(ctx, deviceID, device)
	if err != nil {
		// Check if it's a domain error and convert appropriately
		if domainErr, ok := err.(*errors.DomainError); ok {
			h.logger.Warn("device update failed",
				zap.String("device_id", deviceID),
				zap.String("error_type", string(domainErr.Type)),
				zap.String("operation", domainErr.Operation),
				zap.Error(err),
			)
			return domainErr.ToAPIError().ToResponse(), nil
		}

		// Fallback for unknown errors
		h.logger.Error("unexpected error during device update",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return errors.ErrDeviceUpdateFailed.ToResponse(), nil
	}

	return utils.JSONSuccessResponse(200, updatedDevice), nil
}

func (h *DeviceHandler) CreateDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate and parse request body
	var createReq models.CreateDeviceRequest
	if err := validation.ValidateJSON(request.Body, &createReq); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	// Validate request data
	if err := validation.ValidateCreateDeviceRequest(createReq); err != nil {
		return err.(errors.APIError).ToResponse(), nil
	}

	// Convert to Device model
	device := models.Device{
		MAC:    createReq.MAC,
		Name:   createReq.Name,
		Type:   createReq.Type,
		HomeID: createReq.HomeID,
	}

	h.logger.Debug("creating device",
		zap.String("mac", device.MAC),
		zap.String("name", device.Name),
		zap.String("type", device.Type),
		zap.String("layer", "handler"),
	)

	createdDevice, err := h.svc.CreateDevice(ctx, device)
	if err != nil {
		// Check if it's a domain error and convert appropriately
		if domainErr, ok := err.(*errors.DomainError); ok {
			h.logger.Warn("device creation failed",
				zap.String("device_mac", device.MAC),
				zap.String("error_type", string(domainErr.Type)),
				zap.String("operation", domainErr.Operation),
				zap.Error(err),
			)
			return domainErr.ToAPIError().ToResponse(), nil
		}

		// Fallback for unknown errors
		h.logger.Error("unexpected error during device creation",
			zap.String("device_mac", device.MAC),
			zap.Error(err),
		)
		return errors.ErrDeviceCreationFailed.ToResponse(), nil
	}

	return utils.JSONSuccessResponse(201, createdDevice), nil
}
