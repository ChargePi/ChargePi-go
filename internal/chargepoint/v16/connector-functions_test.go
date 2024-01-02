package v16

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ocppManager-go/ocpp_v16"
)

func newString(s string) *string {
	return &s
}

var (
	ocppConfig = ocpp_v16.Config{
		Version: 1,
		Keys: []core.ConfigurationKey{
			{
				Key:      "AllowOfflineTxForUnknownId",
				Readonly: false,

				Value: newString("false"),
			},
			{
				Key:      "AuthorizationCacheEnabled",
				Readonly: false,
				Value:    newString("false"),
			},
			{
				Key:      "AuthorizeRemoteTxRequests",
				Readonly: false,
				Value:    newString("false"),
			},
			{
				Key:      "ClockAlignedDataInterval",
				Readonly: false,
				Value:    newString("0"),
			},
			{
				Key:      "ConnectionTimeOut",
				Readonly: false,
				Value:    newString("50"),
			},
			{
				Key:      "GetConfigurationMaxKeys",
				Readonly: false,
				Value:    newString("30"),
			},
			{
				Key:      "HeartbeatInterval",
				Readonly: false,
				Value:    newString("60"),
			},
			{
				Key:      "LocalAuthorizeOffline",
				Readonly: false,
				Value:    newString("true"),
			},
			{
				Key:      "LocalPreAuthorize",
				Readonly: false,
				Value:    newString("true"),
			},
			{
				Key:      "MaxEnergyOnInvalidId",
				Readonly: false,
				Value:    newString("0"),
			},
			{
				Key:      "MeterValuesSampledData",
				Readonly: false,
				Value:    newString("Power.Active.Import"),
			},
			{
				Key:      "MeterValuesAlignedData",
				Readonly: false,
				Value:    newString("false"),
			},
			{
				Key:      "NumberOfConnectors",
				Readonly: false,
				Value:    newString("6"),
			},
			{
				Key:      "MeterValueSampleInterval",
				Readonly: false,
				Value:    newString("60"),
			},
			{
				Key:      "ResetRetries",
				Readonly: false,
				Value:    newString("3"),
			},
			{
				Key:      "ConnectorPhaseRotation",
				Readonly: false,
				Value:    newString("0.RST, 1.RST, 2.RTS"),
			},
			{
				Key:      "StopTransactionOnEVSideDisconnect",
				Readonly: false,
				Value:    newString("true"),
			},
			{
				Key:      "StopTransactionOnInvalidId",
				Readonly: false,
				Value:    newString("true"),
			},
			{
				Key:      "StopTxnAlignedData",
				Readonly: false,
			},
			{
				Key:      "StopTxnSampledData",
				Readonly: false,
			},
			{
				Key:      "SupportedFeatureProfiles",
				Readonly: true,
				Value:    newString("Core, LocalAuthListManagement, Reservation, RemoteTrigger"),
			},
			{
				Key:      "TransactionMessageAttempts",
				Readonly: false,
				Value:    newString("3"),
			},
			{
				Key:      "TransactionMessageRetryInterval",
				Readonly: false,
				Value:    newString("60"),
			},
			{
				Key:      "UnlockConnectorOnEVSideDisconnect",
				Readonly: false,
				Value:    newString("true"),
			},
			{
				Key:      "ReserveConnectorZeroSupported",
				Readonly: false,
				Value:    newString("false"),
			},
			{
				Key:      "SendLocalListMaxLength",
				Readonly: false,
				Value:    newString("20"),
			},
			{
				Key:      "LocalAuthListEnabled",
				Readonly: false,
				Value:    newString("true"),
			},
			{
				Key:      "LocalAuthListMaxLength",
				Readonly: false,
				Value:    newString("20"),
			},
		},
	}
)

type connectorFunctionsTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *connectorFunctionsTestSuite) SetupTest() {
	s.cp = &ChargePoint{
		logger: log.StandardLogger(),
	}
}

func (s *connectorFunctionsTestSuite) TestAddConnectors() {
	// todo
}

func (s *connectorFunctionsTestSuite) TestRestoreState() {
	// todo
}

func (s *connectorFunctionsTestSuite) TestNotifyConnectorStatus() {

	s.cp.notifyConnectorStatus(1, core.ChargePointStatusAvailable, core.NoError)

}

func TestConnectorFunctions(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	suite.Run(t, new(connectorFunctionsTestSuite))
}
