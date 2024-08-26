// internal/device/command/create_device.go
package command

import (
	"context"
	"time"

	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"
)

type CreateDeviceHandler struct {
	repo repository.CommandDeviceRepository
}

func NewCreateDeviceHandler(repo repository.CommandDeviceRepository) *CreateDeviceHandler {
	return &CreateDeviceHandler{repo: repo}
}

type CreateDeviceCommand struct {
	Name   string
	Type   string
	UserID string
}

func (h *CreateDeviceHandler) Handle(ctx context.Context, cmd CreateDeviceCommand) (*models.Device, error) {
	device := &models.Device{
		Name:      cmd.Name,
		Type:      cmd.Type,
		UserID:    cmd.UserID,
		State:     "off", // Default state
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return h.repo.CreateDevice(ctx, device)
}
