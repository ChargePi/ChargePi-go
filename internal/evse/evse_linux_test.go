//go:build linux

package evse

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/evcc"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/evcc/mocks"
	powerMeter "github.com/ChargePi/ChargePi-go/pkg/hardware/power-meter"
	powerMeterMock "github.com/ChargePi/ChargePi-go/pkg/hardware/power-meter/mocks"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type evseTestSuite struct {
	suite.Suite
	evse           *V1
	evccMock       *mocks.MockEVCC
	powerMeterMock *powerMeterMock.MockPowerMeter
}

func (s *evseTestSuite) SetupTest() {
	evccMock := mocks.NewMockEVCC(s.T())
	powerMeterMock := powerMeterMock.NewMockPowerMeter(s.T())

	evse, err := NewEvse(1, evccMock, powerMeterMock, 16.0)
	s.Require().NoError(err)

	s.evse = evse
	s.evccMock = evccMock
	s.powerMeterMock = powerMeterMock
}

func (s *evseTestSuite) TestNewEVSEFromSettings() {
	tests := []struct {
		name     string
		settings Settings
		wantErr  bool
	}{
		{
			name: "EVSE settings are valid",
			settings: Settings{
				EvseId:   1,
				MaxPower: 16,
				EVCC: evcc.Settings{
					Type: evcc.TypeDummy,
				},
				PowerMeter: powerMeter.Settings{
					Type: powerMeter.TypeDummy,
				},
				Connectors: nil,
			},
			wantErr: false,
		},
		{
			name: "Missing EVCC",
			settings: Settings{
				EvseId:   1,
				MaxPower: 16,
				PowerMeter: powerMeter.Settings{
					Type: powerMeter.TypeDummy,
				},
				Connectors: nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid EVCC type",
			settings: Settings{
				EvseId:   1,
				MaxPower: 16,
				EVCC: evcc.Settings{
					Type: "unknown",
				},
				PowerMeter: powerMeter.Settings{
					Enabled: true,
					Type:    powerMeter.TypeDummy,
				},
				Connectors: nil,
			},
			wantErr: true,
		},
		{
			name: "PowerMeter disabled",
			settings: Settings{
				EvseId:   1,
				MaxPower: 16,
				EVCC: evcc.Settings{
					Type: evcc.TypeDummy,
				},
				PowerMeter: powerMeter.Settings{
					Enabled: false,
					Type:    powerMeter.TypeDummy,
				},
				Connectors: nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid PowerMeter type",
			settings: Settings{
				EvseId:   1,
				MaxPower: 16,
				EVCC: evcc.Settings{
					Type: evcc.TypeDummy,
				},
				PowerMeter: powerMeter.Settings{
					Enabled: true,
					Type:    "unknown",
				},
				Connectors: nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid max current",
			settings: Settings{
				EvseId:   1,
				MaxPower: -1,
				EVCC: evcc.Settings{
					Type: evcc.TypeDummy,
				},
				PowerMeter: powerMeter.Settings{
					Enabled: false,
					Type:    powerMeter.TypeDummy,
				},
				Connectors: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			_, err := NewEvseFromSettings(tt.settings)
			if tt.wantErr {
				s.Assert().Error(err)
				return
			}

			s.Assert().NoError(err)
		})
	}
}

func (s *evseTestSuite) TestStartCharging() {
	tests := []struct {
		name           string
		connectorId    *int
		measurands     []types.Measurand
		sampleInterval string
		wantErr        bool
	}{
		{
			name: "Started charging",
			measurands: []types.Measurand{
				types.MeasurandCurrentImport,
				types.MeasurandVoltage,
				types.MeasurandPowerActiveImport,
				types.MeasurandEnergyActiveImportRegister,
			},
			sampleInterval: "10s",
			wantErr:        false,
		},
		{
			name: "Unable to enable charging on EVCC",
			measurands: []types.Measurand{
				types.MeasurandCurrentImport,
				types.MeasurandVoltage,
				types.MeasurandPowerActiveImport,
				types.MeasurandEnergyActiveImportRegister,
			},
			sampleInterval: "10s",
			wantErr:        true,
		},
		{
			name: "Power meter is disabled",
			measurands: []types.Measurand{
				types.MeasurandCurrentImport,
				types.MeasurandVoltage,
				types.MeasurandPowerActiveImport,
				types.MeasurandEnergyActiveImportRegister,
			},
			sampleInterval: "10s",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "Unable to enable charging on EVCC" {
				s.evccMock.EXPECT().EnableCharging().Return(errors.New("error starting charging")).Once()
			} else {
				s.evccMock.EXPECT().EnableCharging().Return(nil).Once()
				s.evccMock.EXPECT().Lock().Return().Once()
			}

			err := s.evse.StartCharging(tt.connectorId, tt.measurands, tt.sampleInterval)
			if tt.wantErr {
				s.Assert().Error(err)
				return
			}

			// todo Wait for 10s to allow the goroutine to sample the power meter?

			s.Assert().NoError(err)
		})
	}
}

func (s *evseTestSuite) TestStopCharging() {
	tests := []struct {
		name string
	}{
		{
			name: "Stopped charging",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestSamplePowerMeter() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestCleanup() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *evseTestSuite) TestSetPowerMeter() {
	tests := []struct {
		name       string
		powerMeter powerMeter.PowerMeter
		wantErr    bool
	}{
		{
			name:       "Valid PowerMeter",
			powerMeter: &powerMeter.Dummy{},
			wantErr:    false,
		},
		{
			name:    "PowerMeter nil",
			wantErr: true,
		},
		{
			name:       "PowerMeter disabled",
			powerMeter: &powerMeter.Dummy{},
			wantErr:    true,
		},
		{
			name:       "PowerMeter cleanup failed",
			powerMeter: &powerMeter.Dummy{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "PowerMeter cleanup failed" {
				s.powerMeterMock.EXPECT().Cleanup().Return(errors.New("error")).Once()
			} else if tt.name != "PowerMeter nil" {
				s.powerMeterMock.EXPECT().Cleanup().Return(nil).Once()
			} else if tt.name == "PowerMeter disabled" {
				s.evse.powerMeterEnabled = false
				s.powerMeterMock.EXPECT().Cleanup().Return(nil).Once()
			}

			err := s.evse.SetPowerMeter(tt.powerMeter)
			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
		})
	}
}

func (s *evseTestSuite) TestSetEVCC() {
	tests := []struct {
		name    string
		evcc    evcc.EVCC
		wantErr bool
	}{
		{
			name:    "Valid EVCC",
			evcc:    &evcc.Dummy{},
			wantErr: false,
		},
		{
			name:    "EVCC nil",
			wantErr: true,
		},
		{
			name:    "EVCC cleanup failed",
			evcc:    &evcc.Dummy{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "EVCC cleanup failed" {
				s.evccMock.EXPECT().Cleanup().Return(errors.New("error")).Once()
			} else if tt.name != "EVCC nil" {
				s.evccMock.EXPECT().Cleanup().Return(nil).Once()
			}

			err := s.evse.SetEvcc(tt.evcc)
			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
		})
	}
}

func (s *evseTestSuite) Test_listenForStatusUpdates() {
	tests := []struct {
		name string
		// EVSE internal state
		currentState  state
		expectedState state
		// Notification state
		expectedNotification notifications.StatusNotification
		notification         evcc.StateNotification
		wantErr              bool
	}{
		{
			name: "From available to preparing",
			notification: evcc.StateNotification{
				State: evcc.StateA1,
				Error: "",
			},
			expectedState: state{
				availability: core.AvailabilityTypeOperative,
				status:       core.ChargePointStatusPreparing,
				errorCode:    core.NoError,
			},
			currentState: state{
				availability: core.AvailabilityTypeOperative,
				status:       core.ChargePointStatusAvailable,
				errorCode:    core.NoError,
			},
			expectedNotification: notifications.StatusNotification{
				EvseId:    0,
				Status:    "",
				ErrorCode: "",
			},
			wantErr: false,
		},
		{
			name: "From preparing to charging",
			notification: evcc.StateNotification{
				State: evcc.StateA1,
				Error: "",
			},
			expectedState: state{
				availability: core.AvailabilityTypeOperative,
				status:       core.ChargePointStatusCharging,
				errorCode:    core.NoError,
			},
			currentState: state{
				availability: core.AvailabilityTypeOperative,
				status:       core.ChargePointStatusPreparing,
				errorCode:    core.NoError,
			},
			expectedNotification: notifications.StatusNotification{
				EvseId:    1,
				Status:    "",
				ErrorCode: "",
			},
			wantErr: false,
		},
		{
			name: "From charging to suspendedEV",
			notification: evcc.StateNotification{
				State: evcc.StateA1,
				Error: "",
			},
			expectedState: state{
				availability: "",
				status:       "",
				errorCode:    "",
			},
			currentState: state{
				availability: "",
				status:       "",
				errorCode:    "",
			},
			expectedNotification: notifications.StatusNotification{
				EvseId:    0,
				Status:    "",
				ErrorCode: "",
			},
			wantErr: false,
		},
		{
			name: "From charging to suspendedEVSE",
			notification: evcc.StateNotification{
				State: evcc.StateA1,
				Error: "",
			},
			expectedState: state{
				availability: "",
				status:       "",
				errorCode:    "",
			},
			currentState: state{
				availability: "",
				status:       "",
				errorCode:    "",
			},
			expectedNotification: notifications.StatusNotification{
				EvseId:    0,
				Status:    "",
				ErrorCode: "",
			},
			wantErr: false,
		},
		{
			name: "From suspendedEVSE to charging",
			notification: evcc.StateNotification{
				State: evcc.StateA1,
				Error: "",
			},
			expectedState: state{
				availability: "",
				status:       "",
				errorCode:    "",
			},
			currentState: state{
				availability: "",
				status:       "",
				errorCode:    "",
			},
			expectedNotification: notifications.StatusNotification{
				EvseId:    0,
				Status:    "",
				ErrorCode: "",
			},
			wantErr: false,
		},
		{
			name: "From charging to finishing",
		},
		{
			name: "From charging to faulted",
		},
		{
			name: "From available to faulted",
		},
		{
			name: "From preparing to faulted",
		},
		{
			name: "From available to unavailable",
		},
		{
			name: "From charging to finishing with error",
		},
		{
			name:    "EVCC has no channel",
			wantErr: true,
		},
		{
			name:    "EVCC error",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create channels
			statusNotifications := make(chan notifications.StatusNotification)
			stateChangeChannel := make(chan evcc.StateNotification)

			// Set initial state
			s.evse.SetNotificationChannel(statusNotifications)
			s.evse.state = &tt.currentState

			if tt.name != "EVCC has no channel" {
				s.evccMock.EXPECT().GetStatusChangeChannel().Return(stateChangeChannel).Once()
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					go s.evse.listenForStatusUpdates(ctx)
				})
			} else {
				// Start listening for status updates
				go s.evse.listenForStatusUpdates(ctx)

				// Send state change notification
				go func() {
					time.Sleep(time.Millisecond * 200)
					stateChangeChannel <- tt.notification
				}()
			}

			// Block until the context is done
		Wait:
			for {
				select {
				case <-ctx.Done():
					if tt.wantErr {
						s.Assert().Error(ctx.Err())
					} else {
						s.Assert().NoError(ctx.Err())
					}
					break Wait
				case notification := <-statusNotifications:
					s.Assert().Equal(tt.expectedNotification, notification)
				}
			}

			// Close channels
			close(statusNotifications)
			close(stateChangeChannel)

		})
	}
}

func (s *evseTestSuite) TestAddConnector() {
	tests := []struct {
		name              string
		wantErr           bool
		connector         ConnectorSettings
		expectedConnector ConnectorSettings
	}{
		{
			name: "Valid connector",
			connector: ConnectorSettings{
				ConnectorId: 1,
				Type:        "Type2",
				Status:      "Available",
			},
			expectedConnector: ConnectorSettings{
				ConnectorId: 1,
				Type:        "Type2",
				Status:      "Available",
			},
			wantErr: false,
		},
		{
			name: "Invalid connector",
			connector: ConnectorSettings{
				ConnectorId: -1,
				Type:        "Type2",
				Status:      "Available",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.evse.AddConnector(tt.connector)
			if tt.wantErr {
				s.Assert().Error(err)
				return
			}

			s.Assert().NoError(err)
			s.Assert().Contains(s.evse.connectors, tt.expectedConnector)
			// todo check order?
		})
	}
}

func (s *evseTestSuite) TestGetConnectors() {
	tests := []struct {
		name string
		// Connectors to add to EVSE
		connectors []ConnectorSettings
		// Returned connectors. The difference is the order in which they are returned!
		expectedConnectors []ConnectorSettings
	}{
		{
			name:               "No connectors",
			connectors:         []ConnectorSettings{},
			expectedConnectors: []ConnectorSettings{},
		},
		{
			name: "Connectors present",
			connectors: []ConnectorSettings{
				{
					ConnectorId: 1,
					Type:        "Type2",
					Status:      "Available",
				},
				{
					ConnectorId: 3,
					Type:        "CCS2",
					Status:      "Charging",
				},
				{
					ConnectorId: 2,
					Type:        "Tesla",
					Status:      "Faulted",
				},
			},
			expectedConnectors: []ConnectorSettings{
				{
					ConnectorId: 1,
					Type:        "Type2",
					Status:      "Available",
				},
				{
					ConnectorId: 2,
					Type:        "Tesla",
					Status:      "Faulted",
				},
				{
					ConnectorId: 3,
					Type:        "CCS2",
					Status:      "Charging",
				},
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "Connectors present" {
				for _, c := range tt.connectors {
					err := s.evse.AddConnector(c)
					s.Require().NoError(err)
				}
			}

			connectors := s.evse.GetConnectors()
			s.Assert().EqualValues(tt.expectedConnectors, connectors)
		})
	}
}

func TestEVSE(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(evseTestSuite))
}
