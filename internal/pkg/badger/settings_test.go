package badger

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse"
)

type settingsTestSuite struct {
	suite.Suite
	db     *Database
	tmpDir string
}

func (s *settingsTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
	tempDir, err := os.MkdirTemp("", "badger_settings_test_*")
	s.Require().NoError(err)
	s.tmpDir = tempDir

	s.db, err = NewBadgerDb(tempDir)
	s.Require().NoError(err)
}

func (s *settingsTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	// Remove temporary database file/dir
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *settingsTestSuite) SetupTest() {
	// Clean up the database before each test
	err := s.db.db.DropAll()
	s.Require().NoError(err)
}

func (s *settingsTestSuite) TestSetEvseSettings() {
	tests := []struct {
		name     string
		settings []evse.Settings
	}{
		{
			name: "Set EVSE settings successfully",
			settings: []evse.Settings{
				{
					EvseId:   1,
					MaxPower: 22.0,
					Connectors: []evse.ConnectorSettings{
						{
							ConnectorId: 1,
							Type:        "Type2",
						},
					},
				},
				{
					EvseId:   2,
					MaxPower: 50.0,
					Connectors: []evse.ConnectorSettings{
						{
							ConnectorId: 1,
							Type:        "CCS",
						},
					},
				},
			},
		},
		{
			name:     "Set empty EVSE settings",
			settings: []evse.Settings{},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			err := s.db.SetEvseSettings(context.Background(), tt.settings)

			// Assert
			s.NoError(err)
		})
	}
}

func (s *settingsTestSuite) TestGetEvseSettings() {
	tests := []struct {
		name          string
		setupSettings []evse.Settings
		expectedCount int
	}{
		{
			name: "Get EVSE settings successfully",
			setupSettings: []evse.Settings{
				{
					EvseId:   1,
					MaxPower: 22.0,
					Connectors: []evse.ConnectorSettings{
						{
							ConnectorId: 1,
							Type:        "Type2",
						},
					},
				},
				{
					EvseId:   2,
					MaxPower: 50.0,
					Connectors: []evse.ConnectorSettings{
						{
							ConnectorId: 1,
							Type:        "CCS",
						},
					},
				},
			},
			expectedCount: 2,
		},
		{
			name:          "Get EVSE settings when none exist",
			setupSettings: []evse.Settings{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if len(tt.setupSettings) > 0 {
				err := s.db.SetEvseSettings(nil, tt.setupSettings)
				s.Require().NoError(err)
			}

			// Execute
			settings, err := s.db.GetEvseSettings(nil)

			// Assert
			s.NoError(err)
			s.Len(settings, tt.expectedCount)
		})
	}
}

func (s *settingsTestSuite) TestGetSettings() {
	tests := []struct {
		name          string
		setupSettings *chargepoint.Settings
		expectError   bool
	}{
		{
			name: "Get settings successfully",
			setupSettings: &chargepoint.Settings{
				Info: chargepoint.Info{
					Type:     "AC",
					MaxPower: 22.0,
					FreeMode: false,
				},
				ConnectionSettings: chargepoint.ConnectionSettings{
					Id:              "test-chargepoint",
					ProtocolVersion: "1.6",
					ServerUri:       "ws://localhost:8080",
				},
			},
			expectError: false,
		},
		{
			name:          "Get settings when none exist",
			setupSettings: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupSettings != nil {
				err := s.db.UpdateSettings(context.Background(), *tt.setupSettings)
				s.Require().NoError(err)
			}

			// Execute
			settings, err := s.db.GetSettings(context.Background())

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(settings)
			} else {
				s.NoError(err)
				s.NotNil(settings)
				s.Equal(tt.setupSettings.Info.Type, settings.Info.Type)
				s.Equal(tt.setupSettings.Info.MaxPower, settings.Info.MaxPower)
				s.Equal(tt.setupSettings.ConnectionSettings.Id, settings.ConnectionSettings.Id)
			}
		})
	}
}

func (s *settingsTestSuite) TestUpdateSettings() {
	tests := []struct {
		name        string
		settings    chargepoint.Settings
		expectError bool
	}{
		{
			name: "Update settings successfully",
			settings: chargepoint.Settings{
				Info: chargepoint.Info{
					Type:     "DC",
					MaxPower: 50.0,
					FreeMode: true,
				},
				ConnectionSettings: chargepoint.ConnectionSettings{
					Id:              "updated-chargepoint",
					ProtocolVersion: "2.0.1",
					ServerUri:       "wss://localhost:8080",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			err := s.db.UpdateSettings(context.Background(), tt.settings)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)

				// Verify settings were actually updated
				updatedSettings, err := s.db.GetSettings(context.Background())
				s.NoError(err)
				s.NotNil(updatedSettings)
				s.Equal(tt.settings.Info.Type, updatedSettings.Info.Type)
				s.Equal(tt.settings.Info.MaxPower, updatedSettings.Info.MaxPower)
				s.Equal(tt.settings.ConnectionSettings.Id, updatedSettings.ConnectionSettings.Id)
			}
		})
	}
}

func (s *settingsTestSuite) TestGetOcppConfiguration() {
	tests := []struct {
		name        string
		version     int
		expectError bool
	}{
		{
			name:        "Get OCPP configuration",
			version:     16,
			expectError: true, // Not implemented yet
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			config, err := s.db.GetOcppConfiguration(context.Background(), tt.version)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(config)
			} else {
				s.NoError(err)
				s.NotNil(config)
			}
		})
	}
}

func (s *settingsTestSuite) TestGetLatestOcppConfiguration() {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "Get latest OCPP configuration",
			expectError: true, // Not implemented yet
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			config, err := s.db.GeLatestOcppConfiguration(context.Background())

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(config)
			} else {
				s.NoError(err)
				s.NotNil(config)
			}
		})
	}
}

func TestSettings(t *testing.T) {
	suite.Run(t, new(settingsTestSuite))
}
