package badger

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type tagTestSuite struct {
	suite.Suite
}

func (s *tagTestSuite) SetupTest() {
}

func (s *tagTestSuite) TearDownSuite() {
}

func (s *tagTestSuite) TestAddTagToAuthList() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestRemoveAuthListTag() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestGetLocalAuthListTag() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestGetLocalAuthListTags() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestGetAuthListTagsForVersion() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestRemoveAuthListAllTagsForVersion() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestAddAuthList() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestAddTag() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestRemoveTag() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestGetTag() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestGetTags() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagTestSuite) TestRemoveTags() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestTag(t *testing.T) {
	suite.Run(t, new(tagTestSuite))
}
