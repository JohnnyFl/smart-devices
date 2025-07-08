package errors

import (
	"context"
	"fmt"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeDatabase     ErrorType = "database"
	ErrorTypeExternal     ErrorType = "external"
	ErrorTypeInternal     ErrorType = "internal"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
)

// DomainError represents an error with additional context and metadata
type DomainError struct {
	Type       ErrorType
	Message    string
	Cause      error
	Context    map[string]interface{}
	Operation  string
	Layer      string
	StatusCode int
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to the error
func (e *DomainError) WithContext(key string, value interface{}) *DomainError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithOperation sets the operation that caused the error
func (e *DomainError) WithOperation(operation string) *DomainError {
	e.Operation = operation
	return e
}

// WithLayer sets the layer where the error occurred
func (e *DomainError) WithLayer(layer string) *DomainError {
	e.Layer = layer
	return e
}

// Is implements error comparison for errors.Is()
func (e *DomainError) Is(target error) bool {
	if t, ok := target.(*DomainError); ok {
		return e.Type == t.Type
	}
	return false
}

// NewDomainError creates a new DomainError
func NewDomainError(errorType ErrorType, message string) *DomainError {
	var statusCode int
	switch errorType {
	case ErrorTypeValidation:
		statusCode = 400
	case ErrorTypeNotFound:
		statusCode = 404
	case ErrorTypeConflict:
		statusCode = 409
	case ErrorTypeUnauthorized:
		statusCode = 401
	case ErrorTypeDatabase, ErrorTypeExternal, ErrorTypeInternal:
		statusCode = 500
	default:
		statusCode = 500
	}

	return &DomainError{
		Type:       errorType,
		Message:    message,
		StatusCode: statusCode,
		Context:    make(map[string]interface{}),
	}
}

// WrapError wraps an existing error with domain error context
func WrapError(errorType ErrorType, message string, cause error) *DomainError {
	err := NewDomainError(errorType, message)
	err.Cause = cause
	return err
}

// Predefined domain errors
var (
	// Validation errors
	ErrDomainInvalidDeviceID = NewDomainError(ErrorTypeValidation, "invalid device ID format")
	ErrDomainMissingMAC      = NewDomainError(ErrorTypeValidation, "MAC address is required")
	ErrDomainInvalidMAC      = NewDomainError(ErrorTypeValidation, "invalid MAC address format")
	ErrDomainMissingName     = NewDomainError(ErrorTypeValidation, "device name is required")
	ErrDomainInvalidName     = NewDomainError(ErrorTypeValidation, "device name must be between 1 and 100 characters")
	ErrDomainInvalidType     = NewDomainError(ErrorTypeValidation, "device type must be one of: thermostat, light, camera, sensor")
	ErrDomainInvalidHomeID   = NewDomainError(ErrorTypeValidation, "home ID must be a valid UUID")
	ErrDomainMissingHomeID   = NewDomainError(ErrorTypeValidation, "home ID is required")

	// Not found errors
	ErrDomainDeviceNotFound = NewDomainError(ErrorTypeNotFound, "device not found")
	ErrDomainNoDevicesFound = NewDomainError(ErrorTypeNotFound, "no devices found")

	// Conflict errors
	ErrDomainDeviceExists = NewDomainError(ErrorTypeConflict, "device already exists")

	// Database errors
	ErrDatabaseOperation = NewDomainError(ErrorTypeDatabase, "database operation failed")
	ErrMarshalDevice     = NewDomainError(ErrorTypeDatabase, "failed to marshal device data")
	ErrUnmarshalDevice   = NewDomainError(ErrorTypeDatabase, "failed to unmarshal device data")

	// Internal errors
	ErrInternalOperation = NewDomainError(ErrorTypeInternal, "internal operation failed")
)

// FromContext extracts error context from context.Context if available
func FromContext(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	contextData := make(map[string]interface{})

	// Extract request ID if available
	if requestID := ctx.Value("requestID"); requestID != nil {
		contextData["request_id"] = requestID
	}

	// Extract user ID if available
	if userID := ctx.Value("userID"); userID != nil {
		contextData["user_id"] = userID
	}

	// Extract trace ID if available
	if traceID := ctx.Value("traceID"); traceID != nil {
		contextData["trace_id"] = traceID
	}

	return contextData
}

// ToAPIError converts a DomainError to an APIError for HTTP responses
func (e *DomainError) ToAPIError() APIError {
	var code string

	switch e.Type {
	case ErrorTypeValidation:
		code = "VALIDATION_ERROR"
	case ErrorTypeNotFound:
		code = "NOT_FOUND"
	case ErrorTypeConflict:
		code = "CONFLICT"
	case ErrorTypeUnauthorized:
		code = "UNAUTHORIZED"
	default:
		code = "INTERNAL_ERROR"
	}

	return APIError{
		Code:       code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
	}
}
