package service

import (
	"testing"

	mock_database "github.com/ChargePi/ChargePi-go/gen/mocks/users/pkg/database"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
}

func (s *serviceTestSuite) SetupTest() {
}

func (s *serviceTestSuite) TestGetUsers() {
	dbMock := mock_database.NewMockDatabase(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func (s *serviceTestSuite) TestAddUser() {
	dbMock := mock_database.NewMockDatabase(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func (s *serviceTestSuite) TestGetUser() {
	dbMock := mock_database.NewMockDatabase(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func (s *serviceTestSuite) TestUpdateUser() {
	dbMock := mock_database.NewMockDatabase(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func TestService(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(serviceTestSuite))
}
