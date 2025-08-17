package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	usersv1 "github.com/ChargePi/ChargePi-go/gen/proto/users/v1"
	"github.com/ChargePi/ChargePi-go/internal/users/models"
)

type userTestSuite struct {
	suite.Suite
	server   *grpc.Server
	listener *bufconn.Listener
	// userServiceMock *mock_users.MockUserRepository
	userHandler *UserHandler
}

func (s *userTestSuite) SetupSuite() {
	s.server = grpc.NewServer()
	buffer := 1024 * 1024
	s.listener = bufconn.Listen(buffer)

	// Start the server
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			s.T().Logf("Server failed to serve: %v", err)
		}
	}()
}

func (s *userTestSuite) TearDownSuite() {
	s.server.GracefulStop()
	s.listener.Close()
}

func (s *userTestSuite) SetupTest() {
	// Create mock user service
	// s.userServiceMock = mock_users.NewMockService(s.T())

	// Create user handler
	// s.userHandler = NewUserHandler(s.userServiceMock)
	s.userHandler = &UserHandler{}

	// Register the service
	usersv1.RegisterUserServiceServer(s.server, s.userHandler)
}

func (s *userTestSuite) TestAddUser() {
	tests := []struct {
		name        string
		username    string
		password    string
		role        string
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:        "Add user successfully",
			username:    "testuser",
			password:    "testpass",
			role:        "manufacturer",
			expectError: false,
		},
		{
			name:        "Add user with existing username",
			username:    "existinguser",
			password:    "testpass",
			role:        "technician",
			expectError: true,
			errorCode:   codes.Unknown,
		},
		{
			name:        "Add user with empty username",
			username:    "",
			password:    "testpass",
			role:        "observer",
			expectError: false, // Handler doesn't validate empty username
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup expectations
			if tt.expectError {
				//s.userServiceMock.EXPECT().AddUser(tt.username, tt.password, tt.role).Return(errors.New("user already exists"))
			} else {
				// s.userServiceMock.EXPECT().AddUser(tt.username, tt.password, tt.role).Return(nil)
			}

			// Create request
			request := &usersv1.AddUserRequest{
				User: &usersv1.User{
					Username: tt.username,
					Password: tt.password,
					Role:     tt.role,
				},
			}

			// Execute
			response, err := s.userHandler.AddUser(context.Background(), request)

			// Assert
			if tt.expectError {
				s.Error(err)
				if st, ok := status.FromError(err); ok {
					s.Equal(tt.errorCode, st.Code())
				}
			} else {
				s.NoError(err)
				s.NotNil(response)
				s.Equal("Success", response.Status)
			}
		})
	}
}

func (s *userTestSuite) TestGetUser() {
	tests := []struct {
		name        string
		username    string
		setupUser   *models.User
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:     "Get existing user",
			username: "testuser",
			setupUser: &models.User{
				Username: "testuser",
				Password: "testpass",
				Role:     models.Manufacturer,
			},
			expectError: false,
		},
		{
			name:        "Get non-existing user",
			username:    "nonexistent",
			setupUser:   nil,
			expectError: true,
			errorCode:   codes.Unknown,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup expectations
			if tt.expectError {
				// s.userServiceMock.EXPECT().GetUser(tt.username).Return(nil, errors.New("user doesn't exist"))
			} else {
				// s.userServiceMock.EXPECT().GetUser(tt.username).Return(tt.setupUser, nil)
			}

			// Create request
			request := &usersv1.GetUserRequest{
				Username: tt.username,
			}

			// Execute
			response, err := s.userHandler.GetUser(context.Background(), request)

			// Assert
			if tt.expectError {
				s.Error(err)
				if st, ok := status.FromError(err); ok {
					s.Equal(tt.errorCode, st.Code())
				}
			} else {
				s.NoError(err)
				s.NotNil(response)
				s.NotNil(response.User)
				s.Equal(tt.setupUser.Username, response.User.Username)
				s.Equal(tt.setupUser.Password, response.User.Password)
				s.Equal(string(tt.setupUser.Role), response.User.Role)
			}
		})
	}
}

func (s *userTestSuite) TestGetUsers() {
	tests := []struct {
		name          string
		setupUsers    []models.User
		expectedCount int
	}{
		{
			name: "Get all users",
			setupUsers: []models.User{
				{Username: "user1", Password: "pass1", Role: models.Manufacturer},
				{Username: "user2", Password: "pass2", Role: models.Technician},
				{Username: "user3", Password: "pass3", Role: models.Observer},
			},
			expectedCount: 3,
		},
		{
			name:          "Get users when none exist",
			setupUsers:    []models.User{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup expectations
			// s.userServiceMock.EXPECT().GetUsers().Return(tt.setupUsers, nil)

			// Execute
			response, err := s.userHandler.GetUsers(context.Background(), &usersv1.GetUsersRequest{})

			// Assert
			s.NoError(err)
			s.NotNil(response)
			s.Len(response.Users, tt.expectedCount)
		})
	}
}

func (s *userTestSuite) TestRemoveUser() {
	tests := []struct {
		name        string
		username    string
		setupUser   *models.User
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:     "Remove existing user",
			username: "testuser",
			setupUser: &models.User{
				Username: "testuser",
				Password: "testpass",
				Role:     models.Manufacturer,
			},
			expectError: false,
		},
		{
			name:        "Remove non-existing user",
			username:    "nonexistent",
			setupUser:   nil,
			expectError: true,
			errorCode:   codes.Unknown,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup expectations
			if tt.expectError {
				//	s.userServiceMock.EXPECT().DeleteUser(tt.username).Return(errors.New("user doesn't exist"))
			} else {
				//	s.userServiceMock.EXPECT().DeleteUser(tt.username).Return(nil)
			}

			// Create request
			request := &usersv1.RemoveUserRequest{
				Username: tt.username,
			}

			// Execute
			response, err := s.userHandler.RemoveUser(context.Background(), request)

			// Assert
			if tt.expectError {
				s.Error(err)
				if st, ok := status.FromError(err); ok {
					s.Equal(tt.errorCode, st.Code())
				}
			} else {
				s.NoError(err)
				s.NotNil(response)
				s.Equal("Success", response.Status)
			}
		})
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}
