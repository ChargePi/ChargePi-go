package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/internal/users"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/models"
	"github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
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
	response := &grpc.AddUserResponse{
		Status: "Failed",
	}

	err := s.userService.AddUser(user.GetUsername(), user.GetPassword(), user.GetRole())
	if err == nil {
		response.Status = "Success"
	}

	return response, nil
}

func (s *UserService) GetUser(ctx context.Context, request *grpc.GetUserRequest) (*grpc.User, error) {
	user, err := s.userService.GetUser(request.GetUsername())
	if err != nil {
		return nil, err
	}

	return toUser(*user), nil
}

func (s *UserService) GetUsers(ctx context.Context, e *empty.Empty) (*grpc.GetUsersResponse, error) {
	response := &grpc.GetUsersResponse{}

	getUsers, err := s.userService.GetUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range getUsers {
		response.Users = append(response.Users, toUser(user))
	}

	return response, nil
}

func (s *UserService) RemoveUser(ctx context.Context, request *grpc.RemoveUserRequest) (*grpc.RemoveUserResponse, error) {
	response := &grpc.RemoveUserResponse{
		Status: "Failed",
	}

	err := s.userService.DeleteUser(request.Username)
	if err == nil {
		response.Status = "Success"
	}

	return response, nil
}

func (s *UserService) mustEmbedUnimplementedUsersServer() {
}

func toUser(user models.User) *grpc.User {
	return &grpc.User{
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}
