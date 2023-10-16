package evse

import (
	"testing"
	"time"

	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/pkg/evcc"
	powerMeter "github.com/xBlaz3kx/ChargePi-go/pkg/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type evseTestSuite struct {
	suite.Suite
}

func (s *evseTestSuite) SetupTest() {
}

func (s *evseTestSuite) TestCreateNewEVSE() {
	evccMock := evcc.NewEvccMock(s.T())
	powerMeterMock := powerMeter.NewPowerMeterMock(s.T())

	// Ok case
	connector1, err := NewEvse(1, evccMock, powerMeterMock, 11, nil)
	s.Assert().Equal(1, connector1.evseId)
	s.Assert().Equal(core.ChargePointStatusAvailable, connector1.status)
	s.Assert().Equal(15, connector1.maxChargingTime)
	s.Assert().False(connector1.powerMeterEnabled)

	// Invalid evseId
	_, err = NewEvse(1, evccMock, powerMeterMock, 11, nil)
	s.Assert().Error(err)

	// Invalid evse id
	_, err = NewEvse(0, evccMock, powerMeterMock, 11, nil)
	s.Assert().Error(err)

	// Negative evse id
	_, err = NewEvse(-1, evccMock, powerMeterMock, 11, nil)
	s.Assert().Error(err)
}

func (s *evseTestSuite) TestStartCharging() {
	evccMock := evcc.NewEvccMock(s.T())
	powerMeterMock := powerMeter.NewPowerMeterMock(s.T())

	// Ok case
	evse, err := NewEvse(1, evccMock, powerMeterMock, 11, nil)
	s.Require().NoError(err)

	// Ok case
	err = evse.StartCharging("1234", "exampleTag", nil)
	s.Assert().NoError(err)

	// Cannot start new session on a evse that is already charging
	err = evse.StartCharging("1234a", "exampleTag1", nil)
	s.Assert().Error(err)
	// s.evccMock.AssertNotCalled(s.T(), "Enable")

	err = evse.StopCharging(core.ReasonLocal)
	s.Assert().NoError(err)

	// Invalid transaction and tag id
	err = evse.StartCharging("@1234asd", "~ˇˇ3123", nil)
	s.Assert().Error(err)

	// Invalid transaction id
	err = evse.StartCharging("", "1234", nil)
	s.Assert().Error(err)

	// Invalid tag id
	err = evse.StartCharging("1234", "", nil)
	s.Assert().Error(err)

	// Invalid evse status
	evse.SetStatus(core.ChargePointStatusUnavailable, core.InternalError)
	err = evse.StartCharging("1234a", "1234a", nil)
	s.Assert().Error(err)
}

func (s *evseTestSuite) TestStopCharging() {
	evccMock := evcc.NewEvccMock(s.T())
	powerMeterMock := powerMeter.NewPowerMeterMock(s.T())

	// Ok case
	evse, err := NewEvse(1, evccMock, powerMeterMock, 11, nil)
	s.Require().NoError(err)

	// Start charging
	tagId := util.GenerateRandomTag()
	err = evse.StartCharging("1234", tagId, nil)
	s.Assert().NoError(err)

	// Stop charging normally
	err = evse.StopCharging(core.ReasonLocal)
	s.Assert().NoError(err)

	// Cannot stop charging if the evse is available
	err = evse.StopCharging(core.ReasonLocal)
	s.Assert().Error(err)
}

func (s *evseTestSuite) TestSamplePowerMeter() {
	evccMock := evcc.NewEvccMock(s.T())
	powerMeterMock := powerMeter.NewPowerMeterMock(s.T())

	evse, err := NewEvse(1, evccMock, powerMeterMock, 11, nil)
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