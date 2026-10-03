package v16

import (
	"strings"
	"testing"

	mock_manager "github.com/ChargePi/ChargePi-go/gen/mocks/pkg/configuration/manager"
	"github.com/ChargePi/ocpp-manager/ocpp_v16"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestSetProfilesFromConfigAcceptsRuntimeOnlyProfiles(t *testing.T) {
	profiles := strings.Join([]string{
		core.ProfileName,
		reservation.ProfileName,
		remotetrigger.ProfileName,
		localauth.ProfileName,
	}, ", ")
	settingsManager := mock_manager.NewMockManager(t)

	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.AuthorizationCacheEnabled, mock.Anything).Return(nil)
	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.HeartbeatInterval, mock.Anything).Return(nil)
	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.LightIntensity, mock.Anything).Return(nil)
	settingsManager.EXPECT().GetConfigurationValue(ocpp_v16.SupportedFeatureProfiles).Return(&profiles, nil)
	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.LocalAuthListMaxLength, mock.Anything).Return(nil)

	cp := &ChargePoint{
		chargePoint:     ocpp16.NewChargePoint("test-charge-point", nil, nil),
		settingsManager: settingsManager,
		logger:          zaptest.NewLogger(t),
	}

	require.NoError(t, cp.setProfilesFromConfig())
}

func TestSetProfilesFromConfigRejectsUnknownProfile(t *testing.T) {
	profiles := strings.Join([]string{core.ProfileName, "UnknownProfile"}, ", ")
	settingsManager := mock_manager.NewMockManager(t)

	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.AuthorizationCacheEnabled, mock.Anything).Return(nil)
	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.HeartbeatInterval, mock.Anything).Return(nil)
	settingsManager.EXPECT().OnUpdateKey(ocpp_v16.LightIntensity, mock.Anything).Return(nil)
	settingsManager.EXPECT().GetConfigurationValue(ocpp_v16.SupportedFeatureProfiles).Return(&profiles, nil)

	cp := &ChargePoint{
		chargePoint:     ocpp16.NewChargePoint("test-charge-point", nil, nil),
		settingsManager: settingsManager,
		logger:          zaptest.NewLogger(t),
	}

	require.Error(t, cp.setProfilesFromConfig())
}
