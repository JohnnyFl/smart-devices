package validation

import (
	"encoding/json"
	"regexp"
	"strings"

	"example.com/smart-devices/internal/errors"
	"example.com/smart-devices/internal/models"
	"github.com/google/uuid"
)

var (
	// MAC address regex pattern
	macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
)

// ValidateJSON unmarshals and validates JSON input
func ValidateJSON(body string, target interface{}) error {
	if strings.TrimSpace(body) == "" {
		return errors.ErrMissingRequestBody
	}

	if err := json.Unmarshal([]byte(body), target); err != nil {
		return errors.ErrInvalidJSON.WithDetails("malformed JSON")
	}

	return nil
}

// ValidateDeviceID validates a device ID parameter
func ValidateDeviceID(deviceID string) error {
	if strings.TrimSpace(deviceID) == "" {
		return errors.ErrMissingDeviceID
	}

	// Check if it's a valid UUID format
	if _, err := uuid.Parse(deviceID); err != nil {
		return errors.ErrInvalidRequest.WithMessage("Device ID must be a valid UUID")
	}

	return nil
}

// ValidateCreateDeviceRequest validates a create device request
func ValidateCreateDeviceRequest(req models.CreateDeviceRequest) error {
	var validationErrors []string

	// Validate MAC address
	if req.MAC == "" {
		validationErrors = append(validationErrors, "MAC address is required")
	} else if !macRegex.MatchString(req.MAC) {
		validationErrors = append(validationErrors, "MAC address format is invalid (expected format: XX:XX:XX:XX:XX:XX)")
	}

	// Validate name
	if req.Name == "" {
		validationErrors = append(validationErrors, "name is required")
	} else if len(req.Name) < 1 || len(req.Name) > 100 {
		validationErrors = append(validationErrors, "name must be between 1 and 100 characters")
	}

	// Validate type
	validTypes := map[string]bool{
		"thermostat": true,
		"light":      true,
		"camera":     true,
		"sensor":     true,
	}
	if req.Type == "" {
		validationErrors = append(validationErrors, "type is required")
	} else if !validTypes[req.Type] {
		validationErrors = append(validationErrors, "type must be one of: thermostat, light, camera, sensor")
	}

	// Validate HomeID (UUID format)
	if req.HomeID == "" {
		validationErrors = append(validationErrors, "homeId is required")
	} else if _, err := uuid.Parse(req.HomeID); err != nil {
		validationErrors = append(validationErrors, "homeId must be a valid UUID")
	}

	if len(validationErrors) > 0 {
		return errors.ErrValidationFailed.WithMessage(strings.Join(validationErrors, "; "))
	}

	return nil
}

// ValidateUpdateDeviceRequest validates an update device request
func ValidateUpdateDeviceRequest(req models.UpdateDeviceRequest) error {
	var validationErrors []string

	// Validate name if provided
	if req.Name != nil {
		if len(*req.Name) < 1 || len(*req.Name) > 100 {
			validationErrors = append(validationErrors, "name must be between 1 and 100 characters")
		}
	}

	// Validate type if provided
	if req.Type != nil {
		validTypes := map[string]bool{
			"thermostat": true,
			"light":      true,
			"camera":     true,
			"sensor":     true,
		}
		if !validTypes[*req.Type] {
			validationErrors = append(validationErrors, "type must be one of: thermostat, light, camera, sensor")
		}
	}

	// Validate HomeID if provided (UUID format)
	if req.HomeID != nil {
		if _, err := uuid.Parse(*req.HomeID); err != nil {
			validationErrors = append(validationErrors, "homeId must be a valid UUID")
		}
	}

	// At least one field must be provided for update
	if req.Name == nil && req.Type == nil && req.HomeID == nil {
		validationErrors = append(validationErrors, "at least one field (name, type, or homeId) must be provided for update")
	}

	if len(validationErrors) > 0 {
		return errors.ErrValidationFailed.WithMessage(strings.Join(validationErrors, "; "))
	}

	return nil
}
