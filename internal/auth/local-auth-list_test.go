package auth

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
)

type localAuthListTestSuite struct {
	suite.Suite
	authList LocalAuthList
}

func (s *localAuthListTestSuite) SetupTest() {
	db := database.Get()
	s.authList = NewLocalAuthList(db, 0)
}

func (s *localAuthListTestSuite) TestAddTag() {
	err := s.authList.AddTag("", nil)
	s.Assert().NoError(err)

	err = s.authList.AddTag("", nil)
	s.Assert().NoError(err)

	err = s.authList.AddTag("", nil)
	s.Assert().NoError(err)
}

func (s *localAuthListTestSuite) TestUpdateTag() {
	err := s.authList.UpdateTag("", nil)
	s.Assert().NoError(err)

	err = s.authList.UpdateTag("", nil)
	s.Assert().NoError(err)

	err = s.authList.UpdateTag("", nil)
	s.Assert().NoError(err)
}

func (s *localAuthListTestSuite) TestRemoveTag() {
	err := s.authList.RemoveTag("")
	s.Assert().NoError(err)

	err = s.authList.RemoveTag("")
	s.Assert().NoError(err)

	err = s.authList.RemoveTag("")
	s.Assert().NoError(err)
}

func (s *localAuthListTestSuite) TestRemoveAll() {
	s.authList.RemoveAll()
}

func (s *localAuthListTestSuite) TestGetTag() {
	_, err := s.authList.GetTag("")
	s.Assert().NoError(err)
}

func (s *localAuthListTestSuite) TestGetTags() {
}

func (s *localAuthListTestSuite) TestSetMaxTags() {
}

func (s *localAuthListTestSuite) TestGetVersion() {
}

func (s *localAuthListTestSuite) TestSetVersion() {
}

func TestLocalAuth(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(localAuthListTestSuite))
}
