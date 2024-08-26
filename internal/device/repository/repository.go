package repository

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/pkg/models"
)

type DeviceRepository interface {
	CommandDeviceRepository
	QueryDeviceRepository
}

type CommandDeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error)
	UpdateDeviceState(ctx context.Context, id string, state string) (*models.Device, error)
}

type QueryDeviceRepository interface {
	GetDevice(ctx context.Context, id string) (*models.Device, error)
	ListDevices(ctx context.Context, userID string, page, pageSize int) ([]*models.Device, int, error)
}
