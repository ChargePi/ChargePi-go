package grpc

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/users/mocks"
	grpc2 "github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type userTestSuite struct {
	suite.Suite
	server          *grpc.Server
	listener        *bufconn.Listener
	userServiceMock *mocks.MockUserService
}

func (s *userTestSuite) SetupSuite() {
	s.server = grpc.NewServer()

	buffer := 101024 * 1024
	s.listener = bufconn.Listen(buffer)

	err := s.server.Serve(s.listener)
	s.Require().NoError(err)
}

func (s *userTestSuite) TearDownSuite() {
	// Stop the GRPC server
	err := s.listener.Close()
	s.Require().NoError(err)

	s.server.GracefulStop()
}

func (s *userTestSuite) SetupTest() {
	// Recreate mocks before each test
	s.userServiceMock = mocks.NewMockUserService(s.T())
	grpc2.RegisterUsersServer(s.server, NewUserService(s.userServiceMock))
}

func (s *userTestSuite) TestAddUser() {

}

func (s *userTestSuite) TestUpdateUser() {

}

func (s *userTestSuite) TestRemoveUser() {

}

func (s *userTestSuite) TestGetUser() {

}

func (s *userTestSuite) TestGetUsers() {

}

func TestUser(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(userTestSuite))
}
