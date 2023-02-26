package evse

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type managerTestSuite struct {
	suite.Suite
}

func (s *managerTestSuite) SetupTest() {
}

func (s *managerTestSuite) TestFindEVSE() {
}

func (s *managerTestSuite) TestFindAvailableEVSE() {
}

func (s *managerTestSuite) TestFindEVSEWithReservationId() {
}

func (s *managerTestSuite) TestFindEVSEWithTransactionId() {
}

func (s *managerTestSuite) TestFindEVSEWithTagId() {
}

func (s *managerTestSuite) TestStartCharging() {
}

func (s *managerTestSuite) TestStopCharging() {
}

func (s *managerTestSuite) TestStopAllEVSEs() {
}

func (s *managerTestSuite) TestRestoreEVSEs() {
}

func TestManager(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(managerTestSuite))
}
