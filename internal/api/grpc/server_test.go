package grpc

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type grpcTestSuite struct {
	suite.Suite
}

func (s *grpcTestSuite) SetupTest() {
}

func TestGrpc(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(grpcTestSuite))
}
