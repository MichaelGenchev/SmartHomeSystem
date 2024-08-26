package query

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"
)

type ListDevicesHandler struct {
	repo repository.QueryDeviceRepository
}

func NewListDevicesHandler(repo repository.QueryDeviceRepository) *ListDevicesHandler {
	return &ListDevicesHandler{repo: repo}
}

type ListDevicesQuery struct {
	UserID   string
	Page     int
	PageSize int
}

func (h *ListDevicesHandler) Handle(ctx context.Context, query ListDevicesQuery) ([]*models.Device, int, error) {
	return h.repo.ListDevices(ctx, query.UserID, query.Page, query.PageSize)
}
