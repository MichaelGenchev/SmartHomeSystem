package user

import (
	"context"

	pb "github.com/MichaelGenchev/smart-home-system/pkg/proto"
	"github.com/MichaelGenchev/smart-home-system/internal/user/service"
	"github.com/MichaelGenchev/smart-home-system/internal/user/repository"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"



	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

)

type GRPCServer struct {
	pb.UnimplementedUserServiceServer
	service service.UserService
}

func NewGRPCServer(service service.UserService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.service.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return convertUserToPb(user), nil
}

func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.service.GetUser(ctx, req.Id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return convertUserToPb(user), nil
}

func (s *GRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := s.service.UpdateUser(ctx, req.Id, req.Name, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return convertUserToPb(user), nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.service.DeleteUser(ctx, req.Id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func convertUserToPb(user *models.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}