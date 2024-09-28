package grpc

import (
	"testing"

	grpc2 "github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type tagAuthTestSuite struct {
	suite.Suite
	server   *grpc.Server
	listener *bufconn.Listener
}

func (s *tagAuthTestSuite) SetupSuite() {
	s.server = grpc.NewServer()

	buffer := 101024 * 1024
	s.listener = bufconn.Listen(buffer)

	err := s.server.Serve(s.listener)
	s.Require().NoError(err)
}

func (s *tagAuthTestSuite) TearDownSuite() {
	// Stop the GRPC server
	err := s.listener.Close()
	s.Require().NoError(err)

	s.server.GracefulStop()
}

func (s *tagAuthTestSuite) SetupTest() {
	// Recreate mocks before each test
	authService := NewTagAuthService(nil)
	grpc2.RegisterTagServer(s.server, authService)
}

func (s *tagAuthTestSuite) TestGetAuthorizedCards() {
	tests := []struct {
		name string
	}{
		{
			name: "There are not authorized cards",
		},
		{
			name: "Error retrieving authorized cards",
		},
		{
			name: "There are authorized cards",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagAuthTestSuite) TestAddAuthorizedCards() {
	tests := []struct {
		name string
	}{
		{
			name: "Unable to add authorized card - invalid status",
		},
		{
			name: "Unable to add authorized card - tag already exists",
		},
		{
			name: "Database error while adding authorized card",
		},
		{
			name: "Successfully added authorized card",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagAuthTestSuite) TestRemoveAuthorizedCard() {
	tests := []struct {
		name string
	}{
		{
			name: "Unable to remove authorized card - doesnt exist",
		},
		{
			name: "Unable to remove authorized card - database error",
		},
		{
			name: "Successfully removed authorized card",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestGrpcAuth(t *testing.T) {
	suite.Run(t, new(tagAuthTestSuite))
}
