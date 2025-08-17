package badger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ChargePi/ChargePi-go/internal/users/models"
)

type userTestSuite struct {
	suite.Suite
	db     *Database
	tmpDir string
}

func (s *userTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
	tempDir, err := os.MkdirTemp("", "badger_test_*")
	s.Require().NoError(err)
	s.tmpDir = tempDir

	s.db, err = NewBadgerDb(tempDir)
	s.Require().NoError(err)
}

func (s *userTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	// Remove temporary database file/dir
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *userTestSuite) SetupTest() {
	// Clean up the database before each test
	err := s.db.db.DropAll()
	s.Require().NoError(err)
}

func (s *userTestSuite) TestGetUser() {
	tests := []struct {
		name         string
		username     string
		setupUser    *models.User
		expectError  bool
		expectedUser *models.User
	}{
		{
			name:         "Get existing user",
			username:     "testuser",
			setupUser:    &models.User{Username: "testuser", Password: "password", Role: models.Manufacturer},
			expectError:  false,
			expectedUser: &models.User{Username: "testuser", Password: "password", Role: models.Manufacturer},
		},
		{
			name:         "Get non-existing user",
			username:     "nonexistent",
			setupUser:    nil,
			expectError:  true,
			expectedUser: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupUser != nil {
				err := s.db.AddUser(*tt.setupUser)
				s.Require().NoError(err)
			}

			// Execute
			user, err := s.db.GetUser(tt.username)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(user)
			} else {
				s.NoError(err)
				s.NotNil(user)
				s.Equal(tt.expectedUser.Username, user.Username)
				s.Equal(tt.expectedUser.Password, user.Password)
				s.Equal(tt.expectedUser.Role, user.Role)
			}
		})
	}
}

func (s *userTestSuite) TestGetUsers() {
	tests := []struct {
		name          string
		setupUsers    []models.User
		expectedCount int
	}{
		{
			name: "Get all users",
			setupUsers: []models.User{
				{Username: "user1", Password: "pass1", Role: models.Manufacturer},
				{Username: "user2", Password: "pass2", Role: models.Technician},
				{Username: "user3", Password: "pass3", Role: models.Observer},
			},
			expectedCount: 3,
		},
		{
			name:          "Get users when none exist",
			setupUsers:    []models.User{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			for _, user := range tt.setupUsers {
				err := s.db.AddUser(user)
				s.Require().NoError(err)
			}

			// Execute
			users, err := s.db.GetUsers()

			// Assert
			s.NoError(err)
			s.Len(users, tt.expectedCount)
		})
	}
}

func (s *userTestSuite) TestAddUser() {
	tests := []struct {
		name        string
		user        models.User
		expectError bool
		errorType   error
	}{
		{
			name:        "Add new user successfully",
			user:        models.User{Username: "newuser", Password: "password", Role: models.Manufacturer},
			expectError: false,
		},
		{
			name:        "Add user with existing username",
			user:        models.User{Username: "existinguser", Password: "password", Role: models.Technician},
			expectError: true,
			errorType:   ErrUserExists,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup - add existing user for the second test case
			if tt.errorType == ErrUserExists {
				err := s.db.AddUser(tt.user)
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.AddUser(tt.user)

			// Assert
			if tt.expectError {
				s.Error(err)
				if tt.errorType != nil {
					s.ErrorIs(err, tt.errorType)
				}
			} else {
				s.NoError(err)

				// Verify user was actually added
				addedUser, err := s.db.GetUser(tt.user.Username)
				s.NoError(err)
				s.Equal(tt.user.Username, addedUser.Username)
				s.Equal(tt.user.Password, addedUser.Password)
				s.Equal(tt.user.Role, addedUser.Role)
			}
		})
	}
}

func (s *userTestSuite) TestUpdateUser() {
	tests := []struct {
		name        string
		setupUser   *models.User
		updateUser  models.User
		expectError bool
		errorType   error
	}{
		{
			name:        "Update existing user successfully",
			setupUser:   &models.User{Username: "updateuser", Password: "oldpass", Role: models.Manufacturer},
			updateUser:  models.User{Username: "updateuser", Password: "newpass", Role: models.Technician},
			expectError: false,
		},
		{
			name:        "Update non-existing user",
			setupUser:   nil,
			updateUser:  models.User{Username: "nonexistent", Password: "newpass", Role: models.Observer},
			expectError: true,
			errorType:   ErrUserDoesntExist,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupUser != nil {
				err := s.db.AddUser(*tt.setupUser)
				s.Require().NoError(err)
			}

			// Execute
			user, err := s.db.UpdateUser(tt.updateUser)

			// Assert
			if tt.expectError {
				s.Error(err)
				if tt.errorType != nil {
					s.ErrorIs(err, tt.errorType)
				}
				s.Nil(user)
			} else {
				s.NoError(err)
				s.NotNil(user)

				// Verify user was actually updated
				updatedUser, err := s.db.GetUser(tt.updateUser.Username)
				s.NoError(err)
				s.Equal(tt.updateUser.Username, updatedUser.Username)
				s.Equal(tt.updateUser.Password, updatedUser.Password)
				s.Equal(tt.updateUser.Role, updatedUser.Role)
			}
		})
	}
}

func (s *userTestSuite) TestDeleteUser() {
	tests := []struct {
		name        string
		setupUser   *models.User
		username    string
		expectError bool
	}{
		{
			name:        "Delete existing user successfully",
			setupUser:   &models.User{Username: "deleteuser", Password: "password", Role: models.Manufacturer},
			username:    "deleteuser",
			expectError: false,
		},
		{
			name:        "Delete non-existing user",
			setupUser:   nil,
			username:    "nonexistent",
			expectError: false, // Badger doesn't return error for non-existing keys
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupUser != nil {
				err := s.db.AddUser(*tt.setupUser)
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.DeleteUser(tt.username)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)

				// Verify user was actually deleted
				if tt.setupUser != nil {
					_, err := s.db.GetUser(tt.username)
					s.Error(err) // Should not exist anymore
				}
			}
		})
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}
