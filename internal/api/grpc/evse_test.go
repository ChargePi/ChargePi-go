package grpc

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/evse/manager/mocks"
	grpc2 "github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type evseTestSuite struct {
	suite.Suite
	server      *grpc.Server
	listener    *bufconn.Listener
	evseManager *mocks.MockEVSEManager
}

func (s *evseTestSuite) SetupSuite() {
	s.server = grpc.NewServer()

	buffer := 101024 * 1024
	s.listener = bufconn.Listen(buffer)

	err := s.server.Serve(s.listener)
	s.Require().NoError(err)
}

func (s *evseTestSuite) TearDownSuite() {
	// Stop the GRPC server
	err := s.listener.Close()
	s.Require().NoError(err)

	s.server.GracefulStop()
}

func (s *evseTestSuite) SetupTest() {
	// Recreate mocks and service before each test
	s.evseManager = mocks.NewMockEVSEManager(s.T())
	grpc2.RegisterEvseServer(s.server, NewEvseService(s.evseManager))
}

func (s *evseTestSuite) TestAddEVSE() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestUpdateUser() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestRemoveUser() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestGetEVSE() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestSetEVCC() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestSetPowerMeter() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestGetUsageForEVSE() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestGetEVSEs() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestEVSE(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(evseTestSuite))
}
