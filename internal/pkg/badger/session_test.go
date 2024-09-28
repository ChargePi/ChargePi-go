package badger

// Create a test for the session repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type sessionRepositoryTestSuite struct {
	suite.Suite
	db     *Database
	tmpDir string
}

func (s *sessionRepositoryTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
	tempDir, err := os.CreateTemp("", "test_*")
	s.Require().NoError(err)

	s.tmpDir = tempDir.Name()

	s.db, err = NewBadgerDb(tempDir.Name())
	s.Require().NoError(err)
}

func (s *sessionRepositoryTestSuite) TearDownSuite() {
	// Remove temporary database file/dir
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *sessionRepositoryTestSuite) SetupTest() {

}

func (s *sessionRepositoryTestSuite) TestCreateSession() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionRepositoryTestSuite) TestStopSession() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionRepositoryTestSuite) TestUpdateSession() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestSessionRepository(t *testing.T) {
	suite.Run(t, new(sessionRepositoryTestSuite))
}
