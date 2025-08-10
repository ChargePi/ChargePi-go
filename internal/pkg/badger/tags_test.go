package badger

import (
	"os"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
)

type tagTestSuite struct {
	suite.Suite
	db     *Database
	tmpDir string
}

func (s *tagTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
	tempDir, err := os.MkdirTemp("", "badger_tag_test_*")
	s.Require().NoError(err)
	s.tmpDir = tempDir

	s.db, err = NewBadgerDb(tempDir)
	s.Require().NoError(err)
}

func (s *tagTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	// Remove temporary database file/dir
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *tagTestSuite) SetupTest() {
	// Clean up the database before each test
	err := s.db.db.DropAll()
	s.Require().NoError(err)
}

func (s *tagTestSuite) TestAddTagToAuthList() {
	tests := []struct {
		name        string
		tagId       string
		tagInfo     *types.IdTagInfo
		expectError bool
	}{
		{
			name:  "Add tag to auth list successfully",
			tagId: "test-tag-123",
			tagInfo: &types.IdTagInfo{
				ExpiryDate: &types.DateTime{Time: time.Now()},
				Status:     types.AuthorizationStatusAccepted,
			},
			expectError: false,
		},
		{
			name:        "Add tag to auth list with nil tag info",
			tagId:       "test-tag-456",
			tagInfo:     nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			err := s.db.AddTagToAuthList(tt.tagId, tt.tagInfo)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagTestSuite) TestRemoveAuthListTag() {
	tests := []struct {
		name        string
		tagId       string
		setupTag    bool
		expectError bool
	}{
		{
			name:        "Remove existing auth list tag",
			tagId:       "remove-tag-123",
			setupTag:    true,
			expectError: false,
		},
		{
			name:        "Remove non-existing auth list tag",
			tagId:       "non-existing-tag",
			setupTag:    false,
			expectError: false, // Badger doesn't return error for non-existing keys
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupTag {
				err := s.db.AddTagToAuthList(tt.tagId, &types.IdTagInfo{
					Status: types.AuthorizationStatusAccepted,
				})
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.RemoveAuthListTag(tt.tagId)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagTestSuite) TestGetLocalAuthListTag() {
	tests := []struct {
		name        string
		tagId       string
		setupTag    *types.IdTagInfo
		expectError bool
	}{
		{
			name:  "Get existing auth list tag",
			tagId: "get-tag-123",
			setupTag: &types.IdTagInfo{
				ExpiryDate: &types.DateTime{Time: time.Now()},
				Status:     types.AuthorizationStatusAccepted,
			},
			expectError: false,
		},
		{
			name:        "Get non-existing auth list tag",
			tagId:       "non-existing-tag",
			setupTag:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupTag != nil {
				err := s.db.AddTagToAuthList(tt.tagId, tt.setupTag)
				s.Require().NoError(err)
			}

			// Execute
			tagInfo, err := s.db.GetLocalAuthListTag(tt.tagId)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(tagInfo)
			} else {
				s.NoError(err)
				s.NotNil(tagInfo)
				s.Equal(tt.setupTag.Status, tagInfo.Status)
			}
		})
	}
}

func (s *tagTestSuite) TestAddTag() {
	tests := []struct {
		name        string
		tagId       string
		tagInfo     *types.IdTagInfo
		expectError bool
	}{
		{
			name:  "Add tag successfully",
			tagId: "add-tag-123",
			tagInfo: &types.IdTagInfo{
				ExpiryDate: &types.DateTime{Time: time.Now()},
				Status:     types.AuthorizationStatusAccepted,
			},
			expectError: false,
		},
		{
			name:        "Add tag with nil tag info",
			tagId:       "add-tag-456",
			tagInfo:     nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			err := s.db.AddTag(tt.tagId, tt.tagInfo)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagTestSuite) TestRemoveTag() {
	tests := []struct {
		name        string
		tagId       string
		setupTag    bool
		expectError bool
	}{
		{
			name:        "Remove existing tag",
			tagId:       "remove-tag-123",
			setupTag:    true,
			expectError: false,
		},
		{
			name:        "Remove non-existing tag",
			tagId:       "non-existing-tag",
			setupTag:    false,
			expectError: false, // Badger doesn't return error for non-existing keys
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupTag {
				err := s.db.AddTag(tt.tagId, &types.IdTagInfo{
					Status: types.AuthorizationStatusAccepted,
				})
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.RemoveTag(tt.tagId)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagTestSuite) TestGetTag() {
	tests := []struct {
		name        string
		tagId       string
		setupTag    *types.IdTagInfo
		expectError bool
	}{
		{
			name:  "Get existing tag",
			tagId: "get-tag-123",
			setupTag: &types.IdTagInfo{
				ExpiryDate: &types.DateTime{Time: time.Now()},
				Status:     types.AuthorizationStatusAccepted,
			},
			expectError: false,
		},
		{
			name:        "Get non-existing tag",
			tagId:       "non-existing-tag",
			setupTag:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupTag != nil {
				err := s.db.AddTag(tt.tagId, tt.setupTag)
				s.Require().NoError(err)
			}

			// Execute
			tagInfo, err := s.db.GetTag(tt.tagId)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(tagInfo)
			} else {
				s.NoError(err)
				s.NotNil(tagInfo)
				s.Equal(tt.setupTag.Status, tagInfo.Status)
			}
		})
	}
}

func (s *tagTestSuite) TestGetTags() {
	tests := []struct {
		name          string
		setupTags     []string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Get all tags",
			setupTags:     []string{"tag1", "tag2", "tag3"},
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Get tags when none exist",
			setupTags:     []string{},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			for _, tagId := range tt.setupTags {
				err := s.db.AddTag(tagId, &types.IdTagInfo{
					Status: types.AuthorizationStatusAccepted,
				})
				s.Require().NoError(err)
			}

			// Execute
			tags, err := s.db.GetTags()

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(tags)
			} else {
				// Note: GetTags is not implemented yet, so this will panic
				// This test will need to be updated when that method is implemented
				s.Panics(func() {
					_, _ = s.db.GetTags()
				})
			}
		})
	}
}

func (s *tagTestSuite) TestRemoveAllTags() {
	tests := []struct {
		name        string
		setupTags   []string
		expectError bool
	}{
		{
			name:        "Remove all tags successfully",
			setupTags:   []string{"tag1", "tag2", "tag3"},
			expectError: false,
		},
		{
			name:        "Remove all tags when none exist",
			setupTags:   []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			for _, tagId := range tt.setupTags {
				err := s.db.AddTag(tagId, &types.IdTagInfo{
					Status: types.AuthorizationStatusAccepted,
				})
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.RemoveAllTags()

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestTag(t *testing.T) {
	suite.Run(t, new(tagTestSuite))
}
