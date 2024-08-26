package service

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/internal/device/command"
	"github.com/MichaelGenchev/smart-home-system/internal/device/query"
	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"github.com/MichaelGenchev/smart-home-system/pkg/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error)
	GetDevice(ctx context.Context, id string) (*models.Device, error)
	UpdateDeviceState(ctx context.Context, id string, state string) (*models.Device, error)
	ListDevices(ctx context.Context, userID string, page, pageSize int) ([]*models.Device, int, error)
}

type DeviceService struct {
	proto.UnimplementedDeviceServiceServer
	repo                repository.DeviceRepository
	createDeviceHandler *command.CreateDeviceHandler
	updateDeviceHandler *command.UpdateDeviceStateHandler
	getDeviceHandler    *query.GetDeviceHandler
	listDevicesHandler  *query.ListDevicesHandler
}

func NewDeviceService(repo repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		createDeviceHandler: command.NewCreateDeviceHandler(repo),
		updateDeviceHandler: command.NewUpdateDeviceStateHandler(repo),
		getDeviceHandler:    query.NewGetDeviceHandler(repo),
		listDevicesHandler:  query.NewListDevicesHandler(repo),
	}
}

func (s *DeviceService) CreateDevice(ctx context.Context, req *proto.CreateDeviceRequest) (*proto.Device, error) {
	cmd := command.CreateDeviceCommand{
		Name:   req.Name,
		Type:   req.Type,
		UserID: req.UserId,
	}

	device, err := s.createDeviceHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &proto.Device{
		Id:     device.ID,
		Name:   device.Name,
		Type:   device.Type,
		State:  device.State,
		UserId: device.UserID,
	}, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, req *proto.GetDeviceRequest) (*proto.Device, error) {
	query := query.GetDeviceQuery{
		ID: req.Id,
	}

	device, err := s.getDeviceHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	return &proto.Device{
		Id:     device.ID,
		Name:   device.Name,
		Type:   device.Type,
		State:  device.State,
		UserId: device.UserID,
	}, nil
}

func (s *DeviceService) UpdateDeviceState(ctx context.Context, req *proto.UpdateDeviceStateRequest) (*proto.Device, error) {
	cmd := command.UpdateDeviceStateCommand{
		ID:    req.Id,
		State: req.State,
	}

	device, err := s.updateDeviceHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &proto.Device{
		Id:     device.ID,
		Name:   device.Name,
		Type:   device.Type,
		State:  device.State,
		UserId: device.UserID,
	}, nil
}

func (s *DeviceService) ListDevices(ctx context.Context, req *proto.ListDevicesRequest) (*proto.ListDevicesResponse, error) {
	query := query.ListDevicesQuery{
		UserID:   req.UserId,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	devices, count, err := s.listDevicesHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}
	protoDevices := make([]*proto.Device, len(devices))
	for i, device := range devices {
		protoDevices[i] = convertToProtoDevice(device)
	}
	return &proto.ListDevicesResponse{
		Devices: protoDevices,
		Total:   int32(count),
	}, nil
}

func convertToProtoDevice(device *models.Device) *proto.Device {
	return &proto.Device{
		Id:        device.ID,
		Name:      device.Name,
		Type:      device.Type,
		State:     device.State,
		UserId:    device.UserID,
		CreatedAt: timestamppb.New(device.CreatedAt),
		UpdatedAt: timestamppb.New(device.UpdatedAt),
	}
}
