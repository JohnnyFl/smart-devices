package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/smart-devices/internal/models"
	"go.uber.org/zap"
)

// MockDeviceRepository implements the repository interface for testing
type MockDeviceRepository struct {
	devices map[string]*models.Device
	err     error
}

func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{
		devices: make(map[string]*models.Device),
	}
}

func (m *MockDeviceRepository) GetDevice(_ context.Context, id string) (*models.Device, error) {
	if m.err != nil {
		return nil, m.err
	}
	device, exists := m.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

func (m *MockDeviceRepository) GetDevices(_ context.Context) ([]models.Device, error) {
	if m.err != nil {
		return nil, m.err
	}
	var devices []models.Device
	for _, device := range m.devices {
		devices = append(devices, *device)
	}
	return devices, nil
}

func (m *MockDeviceRepository) CreateDevice(_ context.Context, device models.Device) (models.Device, error) {
	if m.err != nil {
		return device, m.err
	}
	device.ID = "test-id-123"
	device.CreatedAt = time.Now().UnixMilli()
	device.ModifiedAt = device.CreatedAt
	m.devices[device.ID] = &device
	return device, nil
}

func (m *MockDeviceRepository) UpdateDevice(_ context.Context, id string, device models.Device) (*models.Device, error) {
	if m.err != nil {
		return nil, m.err
	}
	existing, exists := m.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}

	// Update fields
	if device.Name != "" {
		existing.Name = device.Name
	}
	if device.Type != "" {
		existing.Type = device.Type
	}
	if device.HomeID != "" {
		existing.HomeID = device.HomeID
	}
	// Ensure ModifiedAt is always greater than the original
	now := time.Now().UnixMilli()
	if now <= existing.ModifiedAt {
		now = existing.ModifiedAt + 1
	}
	existing.ModifiedAt = now

	return existing, nil
}

func (m *MockDeviceRepository) DeleteDevice(_ context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	_, exists := m.devices[id]
	if !exists {
		return errors.New("device not found")
	}
	delete(m.devices, id)
	return nil
}

func (m *MockDeviceRepository) UpdateDeviceHomeID(_ context.Context, id string, homeID string) error {
	if m.err != nil {
		return m.err
	}
	device, exists := m.devices[id]
	if !exists {
		return errors.New("device not found")
	}
	device.HomeID = homeID
	// Ensure ModifiedAt is always greater than the original
	now := time.Now().UnixMilli()
	if now <= device.ModifiedAt {
		now = device.ModifiedAt + 1
	}
	device.ModifiedAt = now
	return nil
}

func (m *MockDeviceRepository) SetError(err error) {
	m.err = err
}

func TestDeviceService_CreateDevice(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockDeviceRepository()
	service := NewDeviceService(mockRepo, logger)

	ctx := context.Background()
	device := models.Device{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "thermostat",
		HomeID: "test-home-id",
	}

	createdDevice, err := service.CreateDevice(ctx, device)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdDevice.ID == "" {
		t.Error("Expected device ID to be set")
	}

	if createdDevice.CreatedAt == 0 {
		t.Error("Expected CreatedAt to be set")
	}

	if createdDevice.ModifiedAt == 0 {
		t.Error("Expected ModifiedAt to be set")
	}
}

func TestDeviceService_GetDevice(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockDeviceRepository()
	service := NewDeviceService(mockRepo, logger)

	ctx := context.Background()

	// First create a device
	device := models.Device{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "thermostat",
		HomeID: "test-home-id",
	}

	createdDevice, _ := service.CreateDevice(ctx, device)

	// Now get the device
	retrievedDevice, err := service.GetDevice(ctx, createdDevice.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedDevice.ID != createdDevice.ID {
		t.Errorf("Expected device ID %s, got %s", createdDevice.ID, retrievedDevice.ID)
	}

	if retrievedDevice.Name != device.Name {
		t.Errorf("Expected device name %s, got %s", device.Name, retrievedDevice.Name)
	}
}

func TestDeviceService_GetDevice_NotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockDeviceRepository()
	service := NewDeviceService(mockRepo, logger)

	ctx := context.Background()

	_, err := service.GetDevice(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent device")
	}
}

func TestDeviceService_UpdateDevice(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockDeviceRepository()
	service := NewDeviceService(mockRepo, logger)

	ctx := context.Background()

	// Create a device first
	device := models.Device{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "thermostat",
		HomeID: "test-home-id",
	}

	createdDevice, _ := service.CreateDevice(ctx, device)

	// Update the device
	updateDevice := models.Device{
		Name: "Updated Device Name",
		Type: "light",
	}

	updatedDevice, err := service.UpdateDevice(ctx, createdDevice.ID, updateDevice)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedDevice.Name != "Updated Device Name" {
		t.Errorf("Expected updated name 'Updated Device Name', got %s", updatedDevice.Name)
	}

	if updatedDevice.Type != "light" {
		t.Errorf("Expected updated type 'light', got %s", updatedDevice.Type)
	}

	if updatedDevice.ModifiedAt <= createdDevice.ModifiedAt {
		t.Error("Expected ModifiedAt to be updated")
	}
}

func TestDeviceService_DeleteDevice(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockDeviceRepository()
	service := NewDeviceService(mockRepo, logger)

	ctx := context.Background()

	// Create a device first
	device := models.Device{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "thermostat",
		HomeID: "test-home-id",
	}

	createdDevice, _ := service.CreateDevice(ctx, device)

	// Delete the device
	err := service.DeleteDevice(ctx, createdDevice.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify device is deleted
	_, err = service.GetDevice(ctx, createdDevice.ID)
	if err == nil {
		t.Error("Expected error when getting deleted device")
	}
}

func TestDeviceService_UpdateDeviceHomeID(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockDeviceRepository()
	service := NewDeviceService(mockRepo, logger)

	ctx := context.Background()

	// Create a device first
	device := models.Device{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "thermostat",
		HomeID: "original-home-id",
	}

	createdDevice, _ := service.CreateDevice(ctx, device)
	originalModifiedAt := createdDevice.ModifiedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update home ID
	newHomeID := "new-home-id"
	err := service.UpdateDeviceHomeID(ctx, createdDevice.ID, newHomeID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the update
	updatedDevice, _ := service.GetDevice(ctx, createdDevice.ID)
	if updatedDevice.HomeID != newHomeID {
		t.Errorf("Expected home ID %s, got %s", newHomeID, updatedDevice.HomeID)
	}

	if updatedDevice.ModifiedAt <= originalModifiedAt {
		t.Error("Expected ModifiedAt to be updated")
	}
}
