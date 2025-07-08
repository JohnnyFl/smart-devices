package handlers

import (
	"context"
	"encoding/json"
	"example.com/smart-devices/internal/models"
	"example.com/smart-devices/internal/services"
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
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Device ID is required"), nil
	}

	h.logger.Debug("fetching device",
		zap.String("device_id", deviceID),
		zap.String("layer", "handler"),
	)

	device, err := h.svc.GetDevice(ctx, deviceID)
	if err != nil {
		h.logger.Warn("device retrieval failed",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)

		return utils.JSONErrorResponse(404, "DEVICE_NOT_FOUND", "Device not found"), nil

	}

	return utils.JSONSuccessResponse(200, device), nil

}

func (h *DeviceHandler) GetDevices(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	devices, err := h.svc.GetDevices(ctx)
	if err != nil {
		h.logger.Warn("device retrieval failed",
			zap.Error(err),
		)

		return utils.JSONErrorResponse(404, "DEVICES_NOT_FOUND", "No devices found"), nil

	}

	return utils.JSONSuccessResponse(200, devices), nil
}

func (h *DeviceHandler) DeleteDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deviceID, ok := request.PathParameters["id"]
	if !ok || deviceID == "" {
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Device ID is required"), nil
	}

	h.logger.Debug("deleting device",
		zap.String("device_id", deviceID),
		zap.String("layer", "handler"),
	)

	err := h.svc.DeleteDevice(ctx, deviceID)
	if err != nil {
		h.logger.Warn("device deletion failed",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)

		return utils.JSONErrorResponse(404, "DEVICE_NOT_DELETED", "Device not deleted"), nil

	}

	return utils.JSONSuccessResponse(200, map[string]string{"message": "Device deleted"}), nil

}

func (h *DeviceHandler) UpdateDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deviceBody := request.Body
	deviceID, ok := request.PathParameters["id"]
	if !ok || deviceID == "" {
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Device ID is required"), nil
	}
	if deviceBody == "" {
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Body is required"), nil
	}

	var device models.Device

	err := json.Unmarshal([]byte(deviceBody), &device)
	if err != nil {
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Invalid device body"), nil
	}

	h.logger.Debug("updating device",
		zap.String("device_id", device.ID),
		zap.String("layer", "handler"),
	)

	updatedDevice, err := h.svc.UpdateDevice(ctx, deviceID, device)
	if err != nil {
		h.logger.Warn("device update failed",
			//zap.String("device_id", updatedDevice.ID),
			zap.Error(err),
		)

		return utils.JSONErrorResponse(404, "DEVICE_NOT_UPDATED", "Device not updated"), nil

	}

	return utils.JSONSuccessResponse(200, updatedDevice), nil

}

func (h *DeviceHandler) CreateDevice(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deviceBody := request.Body
	if deviceBody == "" {
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Device ID is required"), nil
	}

	var device models.Device

	err := json.Unmarshal([]byte(deviceBody), &device)
	if err != nil {
		return utils.JSONErrorResponse(400, "INVALID_REQUEST", "Invalid device body"), nil
	}

	h.logger.Debug("creating device",
		zap.String("layer", "handler"),
	)

	createdDevice, err := h.svc.CreateDevice(ctx, device)
	if err != nil {
		h.logger.Warn("device creation failed",
			zap.Error(err),
		)

		return utils.JSONErrorResponse(404, "DEVICE_CREATION_FAILED", "Device creation failed"), nil

	}

	return utils.JSONSuccessResponse(200, createdDevice), nil

}
