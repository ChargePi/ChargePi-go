package grpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	evsev1 "github.com/ChargePi/ChargePi-go/gen/proto/evse/v1"
)

type evseTestSuite struct {
	suite.Suite
	server      *grpc.Server
	listener    *bufconn.Listener
	evseHandler *EvseHandler
}

func (s *evseTestSuite) SetupSuite() {
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

func (s *evseTestSuite) TearDownSuite() {
	s.server.GracefulStop()
	s.listener.Close()
}

func (s *evseTestSuite) SetupTest() {
	// Create EVSE handler with nil manager for testing
	s.evseHandler = &EvseHandler{}

	// Register the service
	evsev1.RegisterEvseServiceServer(s.server, s.evseHandler)
}

func (s *evseTestSuite) TestGetEVSEs() {
	// Test that the service can handle basic requests
	response, err := s.evseHandler.GetEVSEs(context.Background(), &empty.Empty{})

	// Should not panic even with nil manager
	s.NoError(err)
	s.NotNil(response)
	s.NotNil(response.Evses)
}

func (s *evseTestSuite) TestGetEVSE() {
	// Test that the service can handle basic requests
	request := &evsev1.GetEVSERequest{
		EvseId: 1,
	}

	response, err := s.evseHandler.GetEVSE(context.Background(), request)

	// Should not panic even with nil manager
	s.NoError(err)
	s.NotNil(response)
}

func TestEvse(t *testing.T) {
	suite.Run(t, new(evseTestSuite))
}
