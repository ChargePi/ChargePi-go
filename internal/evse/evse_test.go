package evse

import (
	"github.com/samber/lo"
	"testing"
	"time"

	mock_evcc "github.com/ChargePi/ChargePi-go/gen/mocks/pkg/evcc"
	mock_power_meter "github.com/ChargePi/ChargePi-go/gen/mocks/pkg/power-meter"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/notifications"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type evseTestSuite struct {
	suite.Suite
	evccMock       *mock_evcc.MockEVCC
	powerMeterMock *mock_power_meter.MockPowerMeter
}

func (s *evseTestSuite) SetupTest() {
	s.evccMock = mock_evcc.NewMockEVCC(s.T())
	s.powerMeterMock = mock_power_meter.NewMockPowerMeter(s.T())
}

func (s *evseTestSuite) TestCreateNewEVSE() {

	// Ok case
	connector1, err := NewEvse(1, s.evccMock, s.powerMeterMock, 11, nil)
	s.Assert().Equal(1, connector1.evseId)
	s.Assert().Equal(core.ChargePointStatusAvailable, connector1.status)
	s.Assert().Equal(15, connector1.maxChargingTime)
	s.Assert().False(connector1.powerMeterEnabled)

	// Invalid evseId
	_, err = NewEvse(1, s.evccMock, s.powerMeterMock, 11, nil)
	s.Assert().Error(err)

	// Invalid evse id
	_, err = NewEvse(0, s.evccMock, s.powerMeterMock, 11, nil)
	s.Assert().Error(err)

	// Negative evse id
	_, err = NewEvse(-1, s.evccMock, s.powerMeterMock, 11, nil)
	s.Assert().Error(err)
}

func (s *evseTestSuite) TestStartCharging() {

	// Ok case
	evse, err := NewEvse(1, s.evccMock, s.powerMeterMock, 11, nil)
	s.Require().NoError(err)

	// Ok case
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().NoError(err)

	// Cannot start new session on a evse that is already charging
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().Error(err)
	// s.s.evccMock.AssertNotCalled(s.T(), "Enable")

	err = evse.StopCharging(core.ReasonLocal)
	s.Assert().NoError(err)

	// Invalid transaction and tag id
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().Error(err)

	// Invalid transaction id
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().Error(err)

	// Invalid tag id
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().Error(err)

	// Invalid evse status
	evse.SetStatus(core.ChargePointStatusUnavailable, core.InternalError)
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().Error(err)
}

func (s *evseTestSuite) TestStopCharging() {
	// Ok case
	evse, err := NewEvse(1, s.evccMock, s.powerMeterMock, 11, nil)
	s.Require().NoError(err)

	// Start charging
	err = evse.StartCharging(lo.ToPtr(1), []types.Measurand{}, "60s")
	s.Assert().NoError(err)

	// Stop charging normally
	err = evse.StopCharging(core.ReasonLocal)
	s.Assert().NoError(err)

	// Cannot stop charging if the evse is available
	err = evse.StopCharging(core.ReasonLocal)
	s.Assert().Error(err)
}

func (s *evseTestSuite) TestSamplePowerMeter() {
	evse, err := NewEvse(1, s.evccMock, s.powerMeterMock, 11, nil)
	s.Require().NoError(err)

	var (
		ctx, cancel    = context.WithTimeout(context.Background(), time.Second*30)
		meterValueChan = make(chan notifications.MeterValueNotification)
	)

	defer cancel()
	go func() {
	Loop:
		for {
			select {
			case notif := <-meterValueChan:
				s.Assert().EqualValues(evse.GetEvseId(), notif.EvseId)

				s.Assert().Len(notif.MeterValues, 3)
				s.Assert().EqualValues("1.000", notif.MeterValues[0].SampledValue[0].Value)
				s.Assert().EqualValues(types.MeasurandVoltage, notif.MeterValues[0].SampledValue[0].Measurand)

				s.Assert().EqualValues("1.000", notif.MeterValues[0].SampledValue[0].Value)
				s.Assert().EqualValues(types.MeasurandCurrentImport, notif.MeterValues[0].SampledValue[1].Measurand)

				s.Assert().EqualValues("1.000", notif.MeterValues[0].SampledValue[0].Value)
				s.Assert().EqualValues(types.MeasurandEnergyActiveImportInterval, notif.MeterValues[0].SampledValue[2].Measurand)
			case <-ctx.Done():
				break Loop
			}
		}
	}()

	evse.SetMeterValuesChannel(meterValueChan)
	evse.powerMeterEnabled = true
	evse.SamplePowerMeter([]types.Measurand{types.MeasurandVoltage, types.MeasurandCurrentImport, types.MeasurandEnergyActiveImportInterval})

	time.Sleep(time.Second)

	evse.SamplePowerMeter([]types.Measurand{types.MeasurandVoltage, types.MeasurandCurrentImport, types.MeasurandEnergyActiveImportInterval})
}

func TestEVSE(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(evseTestSuite))
}
