package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/users"
)

type UserService struct {
	grpc.UnimplementedUsersServer
	userService users.Service
}

func NewUserService(userService users.Service) *UserService {
	return &UserService{
		userService: userService,
	}
}

func (s *UserService) AddUser(ctx context.Context, user *grpc.User) (*grpc.AddUserResponse, error) {
	return nil, nil
}

func (s *UserService) GetUser(ctx context.Context, request *grpc.GetUserRequest) (*grpc.User, error) {
	return nil, nil
}

func (s *UserService) GetUsers(ctx context.Context, e *empty.Empty) (*grpc.GetUsersResponse, error) {
	return nil, nil
}

func (s *UserService) RemoveUser(ctx context.Context, request *grpc.RemoveUserRequest) (*grpc.RemoveUserResponse, error) {
	return nil, nil
}

func (s *UserService) mustEmbedUnimplementedUsersServer() {
}
