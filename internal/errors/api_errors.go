package errors

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// APIError represents a standardized API error
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return e.Message
}

// ToResponse converts APIError to Lambda response
func (e APIError) ToResponse() events.APIGatewayProxyResponse {
	body, _ := json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.Code,
		Message: e.Message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: e.StatusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// Predefined API errors
var (
	// 400 Bad Request errors
	ErrInvalidRequest = APIError{
		Code:       "INVALID_REQUEST",
		Message:    "Invalid request format",
		StatusCode: 400,
	}

	ErrMissingDeviceID = APIError{
		Code:       "MISSING_DEVICE_ID",
		Message:    "Device ID is required",
		StatusCode: 400,
	}

	ErrMissingRequestBody = APIError{
		Code:       "MISSING_REQUEST_BODY",
		Message:    "Request body is required",
		StatusCode: 400,
	}

	ErrInvalidJSON = APIError{
		Code:       "INVALID_JSON",
		Message:    "Invalid JSON format in request body",
		StatusCode: 400,
	}

	ErrValidationFailed = APIError{
		Code:       "VALIDATION_FAILED",
		Message:    "Request validation failed",
		StatusCode: 400,
	}

	// 404 Not Found errors
	ErrDeviceNotFound = APIError{
		Code:       "DEVICE_NOT_FOUND",
		Message:    "Device not found",
		StatusCode: 404,
	}

	ErrNoDevicesFound = APIError{
		Code:       "NO_DEVICES_FOUND",
		Message:    "No devices found",
		StatusCode: 404,
	}

	// 500 Internal Server errors
	ErrInternalServer = APIError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An internal server error occurred",
		StatusCode: 500,
	}

	ErrDeviceCreationFailed = APIError{
		Code:       "DEVICE_CREATION_FAILED",
		Message:    "Failed to create device",
		StatusCode: 500,
	}

	ErrDeviceUpdateFailed = APIError{
		Code:       "DEVICE_UPDATE_FAILED",
		Message:    "Failed to update device",
		StatusCode: 500,
	}

	ErrDeviceDeletionFailed = APIError{
		Code:       "DEVICE_DELETION_FAILED",
		Message:    "Failed to delete device",
		StatusCode: 500,
	}
)

// WithMessage creates a new APIError with a custom message
func (e APIError) WithMessage(message string) APIError {
	return APIError{
		Code:       e.Code,
		Message:    message,
		StatusCode: e.StatusCode,
	}
}

// WithDetails creates a new APIError with additional details
func (e APIError) WithDetails(details string) APIError {
	return APIError{
		Code:       e.Code,
		Message:    e.Message + ": " + details,
		StatusCode: e.StatusCode,
	}
}
