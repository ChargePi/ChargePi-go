package grpc

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/diagnostics/mocks"
	grpc2 "github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type loggingTestSuite struct {
	suite.Suite
	server                *grpc.Server
	listener              *bufconn.Listener
	diagnosticServiceMock *mocks.MockDiagnosticsService
}

func (s *loggingTestSuite) SetupSuite() {
	s.server = grpc.NewServer()

	buffer := 101024 * 1024
	s.listener = bufconn.Listen(buffer)

	err := s.server.Serve(s.listener)
	s.Require().NoError(err)
}

func (s *loggingTestSuite) TearDownSuite() {
	// Stop the GRPC server
	err := s.listener.Close()
	s.Require().NoError(err)

	s.server.GracefulStop()
}

func (s *loggingTestSuite) SetupTest() {
	// Recreate mocks before each test
	s.diagnosticServiceMock = mocks.NewMockDiagnosticsService(s.T())
	grpc2.RegisterLogServer(s.server, NewLogService(s.diagnosticServiceMock))
}

func (s *loggingTestSuite) TestGetLogs() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestLogging(t *testing.T) {
	suite.Run(t, new(loggingTestSuite))
}
