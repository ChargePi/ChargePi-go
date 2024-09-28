package badger

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type settingsTestSuite struct {
	suite.Suite
	db *Database
}

func (s *settingsTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
}

func (s *settingsTestSuite) TearDownSuite() {
	// Remove temporary database file/dir
}

func (s *settingsTestSuite) SetupTest() {
}

func (s *settingsTestSuite) TestSetEvseSettings() {

}
func (s *settingsTestSuite) TestGetEvseSettings() {

}

func (s *settingsTestSuite) TestSetSettings() {

}

func (s *settingsTestSuite) TestGetSettings() {

}

func (s *settingsTestSuite) TestUpdateSettings() {

}

func (s *settingsTestSuite) TestGetOcppConfiguration() {

}

func (s *settingsTestSuite) TestGetLatestOcppConfiguration() {

}

func (s *settingsTestSuite) TestSetOcppConfiguration() {

}

func TestSettings(t *testing.T) {
	suite.Run(t, new(settingsTestSuite))
}
