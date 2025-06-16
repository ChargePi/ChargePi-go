package grpc

import (
	"context"

	grpc "github.com/ChargePi/ChargePi-go/gen/proto/users/v1"
	"github.com/ChargePi/ChargePi-go/internal/users/pkg/models"
	"github.com/ChargePi/ChargePi-go/internal/users/service"
	"github.com/golang/protobuf/ptypes/empty"
)

type UserHandler struct {
	grpc.UnimplementedUserServiceServer
	userService service.Service
}

func NewUserHandler(userService service.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (s *UserHandler) AddUser(ctx context.Context, user *grpc.AddUserRequest) (*grpc.AddUserResponse, error) {
	response := &grpc.AddUserResponse{
		Status: "Failed",
	}

	u := user.GetUser()

	err := s.userService.AddUser(u.GetUsername(), u.GetPassword(), u.GetRole())
	if err == nil {
		response.Status = "Success"
	}

	return response, nil
}

func (s *UserHandler) GetUser(ctx context.Context, request *grpc.GetUserRequest) (*grpc.GetUserResponse, error) {
	user, err := s.userService.GetUser(request.GetUsername())
	if err != nil {
		return nil, err
	}

	return &grpc.GetUserResponse{
		User: toUser(*user),
	}, nil
}

func (s *UserHandler) GetUsers(ctx context.Context, e *empty.Empty) (*grpc.GetUsersResponse, error) {
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

func (s *UserHandler) RemoveUser(ctx context.Context, request *grpc.RemoveUserRequest) (*grpc.RemoveUserResponse, error) {
	response := &grpc.RemoveUserResponse{
		Status: "Failed",
	}

	err := s.userService.DeleteUser(request.Username)
	if err == nil {
		response.Status = "Success"
	}

	return response, nil
}

func (s *UserHandler) mustEmbedUnimplementedUsersServer() {
}

func toUser(user models.User) *grpc.User {
	return &grpc.User{
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}
