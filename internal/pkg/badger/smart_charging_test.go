package badger

import (
	"os"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
)

type smartChargingTestSuite struct {
	suite.Suite
	db     *Database
	tmpDir string
}

func (s *smartChargingTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
	tempDir, err := os.MkdirTemp("", "badger_smart_charging_test_*")
	s.Require().NoError(err)
	s.tmpDir = tempDir

	s.db, err = NewBadgerDb(tempDir)
	s.Require().NoError(err)
}

func (s *smartChargingTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	// Remove temporary database file/dir
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *smartChargingTestSuite) SetupTest() {
	// Clean up the database before each test
	err := s.db.db.DropAll()
	s.Require().NoError(err)
}

func (s *smartChargingTestSuite) TestAddProfile() {
	tests := []struct {
		name        string
		profile     *types.ChargingProfile
		expectError bool
	}{
		{
			name: "Add charging profile successfully",
			profile: &types.ChargingProfile{
				ChargingProfileId:      1,
				TransactionId:          123,
				StackLevel:             1,
				ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile,
				ChargingProfileKind:    types.ChargingProfileKindAbsolute,
				RecurrencyKind:         types.RecurrencyKindDaily,
			},
			expectError: false,
		},
		{
			name:        "Add nil charging profile",
			profile:     nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			err := s.db.AddProfile(tt.profile)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *smartChargingTestSuite) TestGetProfile() {
	tests := []struct {
		name         string
		profileId    int
		setupProfile *types.ChargingProfile
		expectError  bool
	}{
		{
			name:      "Get existing profile",
			profileId: 1,
			setupProfile: &types.ChargingProfile{
				ChargingProfileId:      1,
				TransactionId:          123,
				StackLevel:             1,
				ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile,
				ChargingProfileKind:    types.ChargingProfileKindAbsolute,
				RecurrencyKind:         types.RecurrencyKindDaily,
			},
			expectError: false,
		},
		{
			name:         "Get non-existing profile",
			profileId:    999,
			setupProfile: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupProfile != nil {
				err := s.db.AddProfile(tt.setupProfile)
				s.Require().NoError(err)
			}

			// Execute
			profile, err := s.db.GetProfile(tt.profileId)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(profile)
			} else {
				s.NoError(err)
				s.NotNil(profile)
				s.Equal(tt.setupProfile.ChargingProfileId, profile.ChargingProfileId)
				s.Equal(tt.setupProfile.TransactionId, profile.TransactionId)
			}
		})
	}
}

func (s *smartChargingTestSuite) TestRemoveProfile() {
	tests := []struct {
		name         string
		profileId    int
		setupProfile bool
		expectError  bool
	}{
		{
			name:         "Remove existing profile",
			profileId:    1,
			setupProfile: true,
			expectError:  false,
		},
		{
			name:         "Remove non-existing profile",
			profileId:    999,
			setupProfile: false,
			expectError:  false, // Badger doesn't return error for non-existing keys
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupProfile {
				profile := &types.ChargingProfile{
					ChargingProfileId:      tt.profileId,
					TransactionId:          123,
					StackLevel:             1,
					ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile,
					ChargingProfileKind:    types.ChargingProfileKindAbsolute,
					RecurrencyKind:         types.RecurrencyKindDaily,
				}
				err := s.db.AddProfile(profile)
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.RemoveProfile(tt.profileId)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestSmartCharging(t *testing.T) {
	suite.Run(t, new(smartChargingTestSuite))
}
