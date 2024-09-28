package grpc

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type chargePointTestSuite struct {
	suite.Suite
	server   *grpc.Server
	listener *bufconn.Listener
}

func (s *chargePointTestSuite) SetupTest() {
	s.server = grpc.NewServer()

	buffer := 101024 * 1024
	s.listener = bufconn.Listen(buffer)

	err := s.server.Serve(s.listener)
	s.Require().NoError(err)
}

func (s *chargePointTestSuite) TearDownSuite() {
	// Stop the GRPC server
	s.server.GracefulStop()
}

func (s *chargePointTestSuite) TestSetDisplaySettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetDisplaySettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestSetReaderSettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetReaderSettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestSetIndicatorSettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetIndicatorSettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestRestart() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetOCPPVariables() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestSetOCPPVariables() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetOCPPVariable() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetVersion() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestGetStatus() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *chargePointTestSuite) TestChangeChargePointDetails() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestChargePoint(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(chargePointTestSuite))
}
