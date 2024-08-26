package command

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"
)

type UpdateDeviceStateHandler struct {
	repo repository.CommandDeviceRepository
}

func NewUpdateDeviceStateHandler(repo repository.CommandDeviceRepository) *UpdateDeviceStateHandler {
	return &UpdateDeviceStateHandler{repo: repo}
}

type UpdateDeviceStateCommand struct {
	ID    string
	State string
}

func (h *UpdateDeviceStateHandler) Handle(ctx context.Context, cmd UpdateDeviceStateCommand) (*models.Device, error) {
	return h.repo.UpdateDeviceState(ctx, cmd.ID, cmd.State)
}
