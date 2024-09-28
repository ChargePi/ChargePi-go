package manager

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/evse/manager/mocks"
	mocks2 "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager/mocks"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/suite"
)

type managerTestSuite struct {
	suite.Suite
	evseSettingsRepo      *mocks.MockEvseSettingsRepository
	settingsRepo          *mocks2.MockSettingsRepository
	ocppConfigurationRepo *mocks2.MockOcppConfigurationRepository
	ocppManager           ocpp_v16.Manager
}

func (s *managerTestSuite) SetupSuite() {
	cfg, err := ocpp_v16.DefaultConfigurationFromProfiles(core.ProfileName)
	s.Require().NoError(err)

	manager, err := ocpp_v16.NewV16ConfigurationManager(*cfg, core.ProfileName)
	s.Require().NoError(err)
	s.ocppManager = manager
}

func (s *managerTestSuite) TearDownSuite() {}

func (s *managerTestSuite) SetupTest() {
	s.evseSettingsRepo = mocks.NewMockEvseSettingsRepository(s.T())
	s.settingsRepo = mocks2.NewMockSettingsRepository(s.T())
	s.ocppConfigurationRepo = mocks2.NewMockOcppConfigurationRepository(s.T())

	// Replace the configuration
	cfg, err := ocpp_v16.DefaultConfigurationFromProfiles(core.ProfileName)
	s.Require().NoError(err)

	err = s.ocppManager.SetConfiguration(*cfg)
	s.Require().NoError(err)
}

func (s *managerTestSuite) TestGetChargePointSettings() {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
		{
			name: "Error getting settings",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestSetChargePointSettings() {
	tests := []struct {
		name string
	}{
		{
			"Valid settings",
		},
		{
			"Invalid settings",
		},
		{
			"Error saving settings",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestGetOcppConfiguration() {
	tests := []struct {
		name string
	}{
		{},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestSetConfiguration() {
	tests := []struct {
		name string
	}{
		{
			name: "Valid configuration",
		},
		{
			name: "Invalid configuration",
		},
		{
			name: "Unable to update configuration in database",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestUpdateKey() {
	tests := []struct {
		name string
	}{
		{
			name: "Valid key",
		},
		{
			name: "Invalid key",
		},
		{
			name: "Unable to update key in database",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

// Note: No need to test the rest of the methods as they are simple wrappers around the ocppManager

func TestManager(t *testing.T) {
	suite.Run(t, new(managerTestSuite))
}
