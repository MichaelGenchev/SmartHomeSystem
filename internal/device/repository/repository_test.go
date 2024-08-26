
package repository

import (
	"context"
	"testing"

	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the DeviceRepository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	args := m.Called(ctx, device)
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockRepository) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockRepository) UpdateDeviceState(ctx context.Context, id string, state string) (*models.Device, error) {
	args := m.Called(ctx, id, state)
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockRepository) ListDevices(ctx context.Context, userID string, page, pageSize int) ([]*models.Device, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).([]*models.Device), args.Int(1), args.Error(2)
}

func TestCreateDevice(t *testing.T) {
	mockRepo := new(MockRepository)
	device := &models.Device{
		Name:   "Test Device",
		Type:   "Sensor",
		State:  "Off",
		UserID: "user123",
	}

	mockRepo.On("CreateDevice", mock.Anything, device).Return(device, nil)

	createdDevice, err := mockRepo.CreateDevice(context.Background(), device)

	assert.NoError(t, err)
	assert.Equal(t, device, createdDevice)
	mockRepo.AssertExpectations(t)
}

func TestGetDevice(t *testing.T) {
	mockRepo := new(MockRepository)
	device := &models.Device{
		ID:     "device123",
		Name:   "Test Device",
		Type:   "Sensor",
		State:  "Off",
		UserID: "user123",
	}

	mockRepo.On("GetDevice", mock.Anything, "device123").Return(device, nil)

	fetchedDevice, err := mockRepo.GetDevice(context.Background(), "device123")

	assert.NoError(t, err)
	assert.Equal(t, device, fetchedDevice)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDeviceState(t *testing.T) {
	mockRepo := new(MockRepository)
	device := &models.Device{
		ID:     "device123",
		Name:   "Test Device",
		Type:   "Sensor",
		State:  "On",
		UserID: "user123",
	}

	mockRepo.On("UpdateDeviceState", mock.Anything, "device123", "On").Return(device, nil)

	updatedDevice, err := mockRepo.UpdateDeviceState(context.Background(), "device123", "On")

	assert.NoError(t, err)
	assert.Equal(t, "On", updatedDevice.State)
	mockRepo.AssertExpectations(t)
}

func TestListDevices(t *testing.T) {
	mockRepo := new(MockRepository)
	devices := []*models.Device{
		{ID: "device1", Name: "Device 1", Type: "Sensor", State: "Off", UserID: "user123"},
		{ID: "device2", Name: "Device 2", Type: "Actuator", State: "On", UserID: "user123"},
	}

	mockRepo.On("ListDevices", mock.Anything, "user123", 1, 10).Return(devices, 2, nil)

	listedDevices, total, err := mockRepo.ListDevices(context.Background(), "user123", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, devices, listedDevices)
	assert.Equal(t, 2, total)
	mockRepo.AssertExpectations(t)
}