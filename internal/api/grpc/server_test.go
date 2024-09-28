package grpc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type grpcTestSuite struct {
	suite.Suite
}

func (s *grpcTestSuite) SetupTest() {

}

func (s *grpcTestSuite) TestAuthMiddleware() {

}

func (s *grpcTestSuite) TestHealthcheck() {

}

func TestGrpc(t *testing.T) {
	suite.Run(t, new(grpcTestSuite))
}
