package service

import (
	"context"
	"testing"

	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"github.com/MichaelGenchev/smart-home-system/pkg/proto"
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
	service := NewDeviceService(mockRepo)

	req := &proto.CreateDeviceRequest{
		Name:   "Test Device",
		Type:   "Sensor",
		UserId: "user123",
	}

	expectedDevice := &models.Device{
		ID:     "device123",
		Name:   req.Name,
		Type:   req.Type,
		State:  "Off",
		UserID: req.UserId,
	}

	mockRepo.On("CreateDevice", mock.Anything, mock.AnythingOfType("*models.Device")).Return(expectedDevice, nil)

	response, err := service.CreateDevice(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedDevice.ID, response.Id)
	assert.Equal(t, expectedDevice.Name, response.Name)
	assert.Equal(t, expectedDevice.Type, response.Type)
	assert.Equal(t, expectedDevice.State, response.State)
	assert.Equal(t, expectedDevice.UserID, response.UserId)
	mockRepo.AssertExpectations(t)
}

func TestGetDevice(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewDeviceService(mockRepo)

	expectedDevice := &models.Device{
		ID:     "device123",
		Name:   "Test Device",
		Type:   "Sensor",
		State:  "Off",
		UserID: "user123",
	}

	mockRepo.On("GetDevice", mock.Anything, "device123").Return(expectedDevice, nil)

	req := &proto.GetDeviceRequest{Id: "device123"}
	response, err := service.GetDevice(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedDevice.ID, response.Id)
	assert.Equal(t, expectedDevice.Name, response.Name)
	assert.Equal(t, expectedDevice.Type, response.Type)
	assert.Equal(t, expectedDevice.State, response.State)
	assert.Equal(t, expectedDevice.UserID, response.UserId)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDeviceState(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewDeviceService(mockRepo)

	updatedDevice := &models.Device{
		ID:     "device123",
		Name:   "Test Device",
		Type:   "Sensor",
		State:  "On",
		UserID: "user123",
	}

	mockRepo.On("UpdateDeviceState", mock.Anything, "device123", "On").Return(updatedDevice, nil)

	req := &proto.UpdateDeviceStateRequest{Id: "device123", State: "On"}
	response, err := service.UpdateDeviceState(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, updatedDevice.ID, response.Id)
	assert.Equal(t, updatedDevice.Name, response.Name)
	assert.Equal(t, updatedDevice.Type, response.Type)
	assert.Equal(t, "On", response.State)
	assert.Equal(t, updatedDevice.UserID, response.UserId)
	mockRepo.AssertExpectations(t)
}

func TestListDevices(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewDeviceService(mockRepo)

	devices := []*models.Device{
		{ID: "device1", Name: "Device 1", Type: "Sensor", State: "Off", UserID: "user123"},
		{ID: "device2", Name: "Device 2", Type: "Actuator", State: "On", UserID: "user123"},
	}

	mockRepo.On("ListDevices", mock.Anything, "user123", 1, 10).Return(devices, 2, nil)

	req := &proto.ListDevicesRequest{UserId: "user123", Page: 1, PageSize: 10}
	response, err := service.ListDevices(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.Devices))
	assert.Equal(t, int32(2), response.Total)
	for i, device := range response.Devices {
		assert.Equal(t, devices[i].ID, device.Id)
		assert.Equal(t, devices[i].Name, device.Name)
		assert.Equal(t, devices[i].Type, device.Type)
		assert.Equal(t, devices[i].State, device.State)
		assert.Equal(t, devices[i].UserID, device.UserId)
	}
	mockRepo.AssertExpectations(t)
}
