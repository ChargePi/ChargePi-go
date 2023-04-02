package auth

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
)

type localAuthListTestSuite struct {
	suite.Suite
	authList LocalAuthList
}

func (s *localAuthListTestSuite) SetupTest() {
	db := database.Get()
	s.authList = NewLocalAuthList(db, 10)
}

func (s *localAuthListTestSuite) TestAddTag() {
	tagId := util.GenerateRandomTag()
	err := s.authList.AddTag(tagId, okTag)
	s.Assert().NoError(err)

	tagId = util.GenerateRandomTag()
	err = s.authList.AddTag(tagId, blockedTag)
	s.Assert().NoError(err)

	tagId = util.GenerateRandomTag()
	err = s.authList.AddTag(tagId, nil)
	s.Assert().Error(err)

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
	tagId := util.GenerateRandomTag()
	err := s.authList.AddTag(tagId, blockedTag)
	s.Require().NoError(err)

	err = s.authList.RemoveTag(tagId)
	s.Assert().NoError(err)

	err = s.authList.RemoveTag("")
	s.Assert().Error(err)

	tagId = util.GenerateRandomTag()
	err = s.authList.RemoveTag(tagId)
	s.Assert().Error(err)
}

func (s *localAuthListTestSuite) TestRemoveAll() {
	s.authList.RemoveAll()
}

func (s *localAuthListTestSuite) TestGetTag() {
	_, err := s.authList.GetTag("")
	s.Assert().NoError(err)
}

func (s *localAuthListTestSuite) TestGetTags() {
	tagList := s.authList.GetTags()
	s.Assert().NotEmpty(tagList)
}

func (s *localAuthListTestSuite) TestSetMaxTags() {
}

func (s *localAuthListTestSuite) TestVersion() {
	version := s.authList.GetVersion()
	s.Assert().EqualValues(1, version)

	s.authList.SetVersion(1)
	s.Assert().EqualValues(1, version)

	s.authList.SetVersion(0)
	version = s.authList.GetVersion()
	s.Assert().EqualValues(1, version)
}

func TestLocalAuth(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(localAuthListTestSuite))
}
