package users

import (
	"errors"
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/users/mocks"
	"github.com/ChargePi/ChargePi-go/internal/users/models"
	encMock "github.com/ChargePi/ChargePi-go/pkg/encryption/mocks"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	mockRepository *mocks.MockUserRepository
	mockEncryptor  *encMock.MockEncryptor
	service        *UserService
}

func (s *serviceTestSuite) SetupTest() {
	s.mockRepository = mocks.NewMockUserRepository(s.T())
	s.service = NewUserService(s.mockRepository)
}

func (s *serviceTestSuite) TestGetUsers() {
	tests := []struct {
		name          string
		expectedUsers []*models.User
		err           error
	}{
		{
			name: "No users",
		},
		{
			name: "One user",
		},
		{
			name: "Multiple users",
		},
		{
			name: "Database error",
		},
		{
			name: "Unauthorized",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			user, err := s.service.GetUsers()
			s.Assert().NoError(err)
			s.Assert().Equal(0, len(user))
		})
	}
}

func (s *serviceTestSuite) TestAddUser() {
	tests := []struct {
		name     string
		username string
		password string
		role     models.Role
		err      error
	}{
		{
			name: "Valid user",
		},
		{
			name: "User already exists",
		},
		{
			name: "Encryption failed",
		},
		{
			name: "Validation failed - no password",
		},
		{
			name: "Validation failed - no username",
		},
		{
			name: "Validation failed - invalid role",
		},
		{
			name: "Database error",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.service.AddUser(tt.username, tt.password, string(tt.role))
			s.Assert().NoError(err)
		})
	}
}

func (s *serviceTestSuite) TestGetUser() {
	validUser := &models.User{
		Username: "exampleUser",
		Password: "examplePassword",
		Role:     models.Manufacturer,
	}
	s.mockRepository.EXPECT().GetUser("exampleUser").Return(validUser, nil)
	s.mockRepository.EXPECT().GetUser("notFound").Return(nil, errors.New("user not found"))
	s.mockRepository.EXPECT().GetUser("exampleUser1").Return(nil, errors.New("database error"))

	tests := []struct {
		name         string
		username     string
		expectedUser *models.User
		err          error
	}{
		{
			name: "User found",
		},
		{
			name: "User not found",
		},
		{
			name: "Other database error",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.service.GetUser("exampleUser")
			if tt.err != nil {
				s.Assert().Error(err)
				s.Assert().ErrorIs(tt.err, err)
				s.Assert().Nil(user)
				return
			}

			s.Assert().NoError(err)
			s.Assert().Equal(tt.expectedUser, user)
		})
	}
}

func (s *serviceTestSuite) TestUpdateUser() {
	s.mockRepository.EXPECT().UpdateUser(gomock.Any()).Return(nil, nil).Times(1)

	tests := []struct {
		name         string
		expectedUser *models.User
		username     string
		password     *string
		role         *string
		err          error
	}{
		{
			name: "Valid user",
		},
		{
			name: "User not found",
		},
		{
			name: "Invalid role",
		},
		{
			name: "Unauthorized",
		},
		{
			name: "Database error",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// User exists
			user, err := s.service.UpdateUser("", lo.ToPtr(""), lo.ToPtr(""))
			if tt.err != nil {
				s.Assert().Error(err)
				s.Assert().ErrorIs(tt.err, err)
				s.Assert().Nil(user)
				return
			}

			s.Assert().NoError(err)
			s.Assert().Equal("", user.Username)
		})
	}
}

func TestService(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(serviceTestSuite))
}
