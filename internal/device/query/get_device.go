package query

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"
)

type GetDeviceHandler struct {
	repo repository.QueryDeviceRepository
}

func NewGetDeviceHandler(repo repository.QueryDeviceRepository) *GetDeviceHandler {
	return &GetDeviceHandler{repo: repo}
}

func (h *GetDeviceHandler) Handle(ctx context.Context, query GetDeviceQuery) (*models.Device, error) {
	return h.repo.GetDevice(ctx, query.ID)
}

type GetDeviceQuery struct {
	ID string
}
