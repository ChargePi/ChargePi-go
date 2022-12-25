package session

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SessionTestSuite struct {
	suite.Suite
	emptySession Session
	validSession Session
}

func (s *SessionTestSuite) SetupTest() {
	s.emptySession = Session{
		IsActive:      false,
		TransactionId: "",
		Started:       "",
		Consumption:   nil,
	}

	s.validSession = Session{
		IsActive:      true,
		TransactionId: "1234",
		TagId:         "1234",
		Started:       time.Now().Format(time.RFC3339),
		Consumption:   nil,
	}
}

func (s *SessionTestSuite) TestAddSampledValue() {
	var (
		samples = []types.SampledValue{
			{
				Value:     "123.21",
				Measurand: types.MeasurandCurrentImport,
			}, {
				Value:     "123.21",
				Measurand: types.MeasurandVoltage,
			}, {
				Value:     "123.21",
				Measurand: types.MeasurandPowerActiveImport,
			},
		}

		expected = []types.MeterValue{
			{
				SampledValue: samples,
			},
		}
	)

	// Session not active
	s.emptySession.AddSampledValue(samples)
	s.Require().Nil(s.emptySession.Consumption)

	// Start session
	hasStarted := s.emptySession.StartSession("1234", "1234")
	s.Require().NoError(hasStarted)

	// Session active
	s.emptySession.AddSampledValue(samples)
	s.Require().EqualValues(expected, s.emptySession.Consumption)
}

func (s *SessionTestSuite) TestStartSession() {
	var expectedSession = Session{
		IsActive:      true,
		TransactionId: "test1234",
		TagId:         "test1234",
		Started:       time.Now().Format(time.RFC3339),
		Consumption:   []types.MeterValue{},
	}

	err := s.emptySession.StartSession("", "")
	s.Require().Error(err)
	s.Require().EqualValues(s.emptySession, s.emptySession)

	err = s.emptySession.StartSession("1234", "")
	s.Require().Error(err)
	s.Require().EqualValues(s.emptySession, s.emptySession)

	err = s.emptySession.StartSession("", "test1234")
	s.Require().Error(err)
	s.Require().EqualValues(s.emptySession, s.emptySession)

	// Ok case
	err = s.emptySession.StartSession("test1234", "test1234")
	s.Require().NoError(err)
	s.Require().EqualValues(expectedSession, s.emptySession)
}

func (s *SessionTestSuite) TestEndSession() {
	s.validSession.EndSession()
	s.Require().EqualValues(s.validSession, s.emptySession)
}

func (s *SessionTestSuite) TestCalculateAvgPower() {
	var (
		emptySample = []types.SampledValue{}
		zeroSamples = []types.MeterValue{
			{
				Timestamp:    &types.DateTime{Time: time.Now()},
				SampledValue: emptySample,
			}, {
				Timestamp:    &types.DateTime{Time: time.Now()},
				SampledValue: emptySample,
			},
		}
		meterWithSamplePower10 = types.MeterValue{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "10",
					Measurand: types.MeasurandPowerActiveImport,
				},
			},
		}

		consumption1 = []types.MeterValue{
			meterWithSamplePower10, meterWithSamplePower10,
		}

		mixedConsumption = []types.MeterValue{
			meterWithSamplePower10,
			{
				Timestamp: &types.DateTime{Time: time.Now()},
				SampledValue: []types.SampledValue{
					{
						Value:     "2",
						Measurand: types.MeasurandCurrentImport,
					}, {
						Value:     "15",
						Measurand: types.MeasurandVoltage,
					},
				},
			},
		}

		faultyMixedConsumption = []types.MeterValue{
			meterWithSamplePower10,
			{
				Timestamp: &types.DateTime{Time: time.Now()},
				SampledValue: []types.SampledValue{
					{
						Value:     "2",
						Measurand: types.MeasurandCurrentImport,
					},
				},
			},
		}

		threeMeasurands = []types.MeterValue{
			{
				Timestamp: &types.DateTime{Time: time.Now()},
				SampledValue: []types.SampledValue{
					{
						Value:     "10",
						Measurand: types.MeasurandPowerActiveImport,
					},
					{
						Value:     "2",
						Measurand: types.MeasurandCurrentImport,
					},
					{
						Value:     "15",
						Measurand: types.MeasurandVoltage,
					},
				},
			}, {
				Timestamp: &types.DateTime{Time: time.Now()},
				SampledValue: []types.SampledValue{
					{
						Value:     "2",
						Measurand: types.MeasurandCurrentImport,
					},
					{
						Value:     "15",
						Measurand: types.MeasurandVoltage,
					},
					{
						Value:     "10",
						Measurand: types.MeasurandPowerActiveImport,
					},
				},
			},
		}
	)

	// Start with zero samples
	s.emptySession.Consumption = zeroSamples
	s.Require().InDelta(0, s.emptySession.CalculateAvgPower(), 0.1)

	s.emptySession.Consumption = consumption1
	s.Require().InDelta(10, s.emptySession.CalculateAvgPower(), 1)

	s.emptySession.Consumption = mixedConsumption
	s.Require().InDelta(20, s.emptySession.CalculateAvgPower(), .5)

	s.emptySession.Consumption = faultyMixedConsumption
	s.Require().InDelta(10, s.emptySession.CalculateAvgPower(), 0.1)

	s.emptySession.Consumption = threeMeasurands
	s.Require().InDelta(10, s.emptySession.CalculateAvgPower(), 0.1)

}

func (s *SessionTestSuite) TestCalculateEnergyConsumptionWithAvgPower() {
	var (
		started5min = time.Now().Add(-5 * time.Minute).Format(time.RFC3339)
	)

	// Start the session
	err := s.emptySession.StartSession("1234", "1234")
	s.Require().NoError(err)

	// Add a sample
	s.emptySession.AddSampledValue([]types.SampledValue{
		{
			Value:     "10",
			Measurand: types.MeasurandPowerActiveImport,
		},
	})

	// Session started 5 minutes before
	s.emptySession.Started = started5min
	s.Require().InDelta(float64(300*10), s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 8)

	// Session just started
	s.emptySession.Started = time.Now().Format(time.RFC3339)
	s.Require().InDelta(0.0, s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 1)

	// Session ended just now
	s.emptySession.EndSession()
	s.Require().InDelta(0.0, s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 1)

	s.emptySession.AddSampledValue([]types.SampledValue{
		{
			Value:     "10",
			Measurand: types.MeasurandPowerActiveImport,
		},
	})
	s.Require().InDelta(0.0, s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 1)
}

func TestSession(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
