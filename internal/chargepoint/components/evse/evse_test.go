package evse

import (
	"os/exec"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/session"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"golang.org/x/net/context"
)

const (
	fileName = "./evse-1.json"
)

type (
	EvseTestSuite struct {
		suite.Suite
		evse           *Impl
		evccMock       *test.EvccMock
		powerMeterMock *test.PowerMeterMock
	}
)

/*---------------------- Test suite ----------------------*/

func NewConnectorTestSuite() *EvseTestSuite {
	return &EvseTestSuite{}
}

func (s *EvseTestSuite) SetupTest() {
	cmd := exec.Command("touch", fileName)
	err := cmd.Run()
	s.Require().NoError(err)

	s.evccMock = new(test.EvccMock)
	s.powerMeterMock = new(test.PowerMeterMock)

	s.evccMock.On("EnableCharging").Return()
	s.evccMock.On("DisableCharging").Return()

	// Create a new evse
	evse, err := NewEvse(
		1,
		s.evccMock,
		s.powerMeterMock,
		false,
		11,
		nil,
	)
	s.Require().NoError(err)

	s.evse = evse
}

func (s *EvseTestSuite) TestCreateNewConnector() {
	s.evccMock.On("DisableCharging").Return()
	s.evccMock.On("EnableCharging").Return()

	// Ok case
	connector1, err := NewEvse(1, s.evccMock, s.powerMeterMock, false, 11, nil)

	s.Require().Equal(1, connector1.evseId)
	s.Require().Equal(core.ChargePointStatusAvailable, connector1.status)
	s.Require().Equal(15, connector1.maxChargingTime)
	s.Require().False(connector1.powerMeterEnabled)

	// Invalid evseId
	_, err = NewEvse(1, s.evccMock, s.powerMeterMock, false, 11, nil)
	s.Require().Error(err)

	// Invalid evse id
	_, err = NewEvse(0, s.evccMock, s.powerMeterMock, false, 11, nil)
	s.Require().Error(err)

	// Negative evse id
	_, err = NewEvse(-1, s.evccMock, s.powerMeterMock, false, 11, nil)
	s.Require().Error(err)
}

func (s *EvseTestSuite) TestStartCharging() {
	// Ok case
	err := s.evse.StartCharging("1234", "exampleTag", nil)
	s.Require().NoError(err)
	s.evccMock.AssertCalled(s.T(), "Enable")

	// Cannot start new session on a evse that is already charging
	err = s.evse.StartCharging("1234a", "exampleTag1", nil)
	s.Require().Error(err)
	// s.evccMock.AssertNotCalled(s.T(), "Enable")

	err = s.evse.StopCharging(core.ReasonLocal)
	s.Require().NoError(err)

	// Invalid transaction and tag id
	err = s.evse.StartCharging("@1234asd", "~ˇˇ3123", nil)
	s.Require().Error(err)
	// s.evccMock.AssertNotCalled(s.T(), "Enable")

	// Invalid transaction id
	err = s.evse.StartCharging("", "1234", nil)
	s.Require().Error(err)
	// s.evccMock.AssertNotCalled(s.T(), "Enable")

	// Invalid tag id
	err = s.evse.StartCharging("1234", "", nil)
	s.Require().Error(err)
	// s.evccMock.AssertNotCalled(s.T(), "Enable")

	// Invalid evse status
	s.evse.SetStatus(core.ChargePointStatusUnavailable, core.InternalError)
	err = s.evse.StartCharging("1234a", "1234a", nil)
	s.Require().Error(err)
}

func (s *EvseTestSuite) TestStopCharging() {
	// Start charging
	err := s.evse.StartCharging("1234", "1234", nil)
	s.Require().NoError(err)
	s.evccMock.AssertCalled(s.T(), "Enable")

	// Stop charging normally
	err = s.evse.StopCharging(core.ReasonLocal)
	s.Require().NoError(err)
	s.evccMock.AssertCalled(s.T(), "Disable")

	// Cannot stop charging if the evse is available
	err = s.evse.StopCharging(core.ReasonLocal)
	s.Require().Error(err)
	// s.evccMock.AssertNotCalled(s.T(), "Disable")
}

func (s *EvseTestSuite) TestResumeCharging() {
	var (
		maxChargingTime = s.evse.GetMaxChargingTime()
		validSession    = session.Session{
			IsActive:      true,
			TransactionId: "1234",
			TagId:         "1234",
			Started:       time.Now().Format(time.RFC3339),
			Consumption:   nil,
		}

		expiredSession = session.Session{
			IsActive:      true,
			TransactionId: "1234",
			TagId:         "1234",
			Started:       time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local).Format(time.RFC3339),
			Consumption:   nil,
		}

		invalidSession = session.Session{
			IsActive:      false,
			TransactionId: "",
			TagId:         "",
			Started:       "",
			Consumption:   nil,
		}
	)

	// Ok case
	s.evse.SetStatus(core.ChargePointStatusCharging, core.NoError)
	timeElapsed, err := s.evse.ResumeCharging(validSession)
	s.Require().NoError(err)
	s.Require().InDelta(0, timeElapsed, 1)

	err = s.evse.StopCharging(core.ReasonLocal)
	s.Require().NoError(err)

	// Invalid session
	s.evse.SetStatus(core.ChargePointStatusCharging, core.NoError)
	timeElapsed, err = s.evse.ResumeCharging(invalidSession)
	s.Require().Error(err)
	s.Require().InDelta(maxChargingTime, timeElapsed, 1)

	// Invalid evse status
	s.evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	timeElapsed, err = s.evse.ResumeCharging(validSession)
	s.Require().Error(err)
	s.Require().InDelta(maxChargingTime, timeElapsed, 1)

	// Expired session
	s.evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	timeElapsed, err = s.evse.ResumeCharging(expiredSession)
	s.Require().Error(err)
	s.Require().InDelta(maxChargingTime, timeElapsed, 1)

	// Invalid status
	s.evse.SetStatus(core.ChargePointStatusUnavailable, core.EVCommunicationError)
	timeElapsed, err = s.evse.ResumeCharging(invalidSession)
	s.Require().Error(err)
	s.Require().InDelta(maxChargingTime, timeElapsed, 1)
}

func (s *EvseTestSuite) TestReserveConnector() {
	// Ok case
	err := s.evse.Reserve(1, "")
	s.Require().NoError(err)

	// EVSE already reserved
	err = s.evse.Reserve(2, "")
	s.Require().Error(err)

	err = s.evse.RemoveReservation()
	s.Require().NoError(err)

	// Invalid evse status
	s.evse.SetStatus(core.ChargePointStatusCharging, core.NoError)
	err = s.evse.Reserve(2, "")
	s.Require().Error(err)
}

func (s *EvseTestSuite) TestRemoveReservation() {
	// Make a reservation
	err := s.evse.Reserve(1, "")
	s.Require().NoError(err)

	// Ok case
	err = s.evse.RemoveReservation()
	s.Require().NoError(err)

	// Cannot remove reservation that is not there
	err = s.evse.RemoveReservation()
	s.Require().Error(err)
}

func (s *EvseTestSuite) TestSamplePowerMeter() {
	s.powerMeterMock.On("GetEnergy").Return(1.0)
	s.powerMeterMock.On("GetCurrent").Return(1.0)
	s.powerMeterMock.On("GetVoltage").Return(1.0)

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
				s.Assert().EqualValues(s.evse.GetEvseId(), notif.EvseId)

				s.Assert().Len(notif.MeterValues, 3)
				s.Assert().EqualValues("1.000", notif.MeterValues[0].SampledValue[0].Value)
				s.Assert().EqualValues(types.MeasurandVoltage, notif.MeterValues[0].SampledValue[0].Measurand)

				s.Assert().EqualValues("1.000", notif.MeterValues[1].SampledValue[0].Value)
				s.Assert().EqualValues(types.MeasurandCurrentImport, notif.MeterValues[1].SampledValue[0].Measurand)

				s.Assert().EqualValues("1.000", notif.MeterValues[2].SampledValue[0].Value)
				s.Assert().EqualValues(types.MeasurandEnergyActiveImportInterval, notif.MeterValues[2].SampledValue[0].Measurand)
			case <-ctx.Done():
				break Loop
			}
		}
	}()

	s.evse.SetMeterValuesChannel(meterValueChan)
	s.evse.powerMeterEnabled = true
	s.evse.powerMeter = s.powerMeterMock
	s.evse.SamplePowerMeter([]types.Measurand{types.MeasurandVoltage, types.MeasurandCurrentImport, types.MeasurandEnergyActiveImportInterval})

	time.Sleep(time.Second)

	s.evse.SamplePowerMeter([]types.Measurand{types.MeasurandVoltage, types.MeasurandCurrentImport, types.MeasurandEnergyActiveImportInterval})
}

func TestConnector(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, NewConnectorTestSuite())
}
