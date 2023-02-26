package service

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/database"
)

type serviceTestSuite struct {
	suite.Suite
}

func (s *serviceTestSuite) SetupTest() {
}

func (s *serviceTestSuite) TestGetUsers() {
	dbMock := database.NewDatabaseMock(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func (s *serviceTestSuite) TestAddUser() {
	dbMock := database.NewDatabaseMock(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func (s *serviceTestSuite) TestGetUser() {
	dbMock := database.NewDatabaseMock(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func (s *serviceTestSuite) TestUpdateUser() {
	dbMock := database.NewDatabaseMock(s.T())
	service := NewUserService(dbMock)

	user, err := service.GetUser("")
	s.Assert().NoError(err)
	s.Assert().Equal("", user.Username)
}

func TestService(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(serviceTestSuite))
}
