package badger

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type userTestSuite struct {
	suite.Suite
	db *Database
}

func (s *userTestSuite) SetupSuite() {
}

func (s *userTestSuite) SetupTest() {
}

func (s *userTestSuite) TearDownSuite() {
}

func (s *userTestSuite) TestGetUser() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *userTestSuite) TestGetUsers() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *userTestSuite) TestAddUser() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *userTestSuite) TestUpdateUser() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *userTestSuite) TestDeleteUser() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}
