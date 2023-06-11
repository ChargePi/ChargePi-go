package session

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
)

type sessionTestSuite struct {
	suite.Suite
	emptySession Session
	validSession Session
}

func (s *sessionTestSuite) SetupTest() {
	currentTime := time.Now()

	s.validSession = Session{
		IsActive:      true,
		TransactionId: "1234",
		TagId:         "1234",
		Started:       &currentTime,
	}
}

func (s *sessionTestSuite) TestAddSampledValue() {
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
	s.Assert().Nil(s.emptySession.Consumption)

	// Start session
	hasStarted := s.emptySession.StartSession("1234", "1234")
	s.Assert().NoError(hasStarted)

	// Session active
	s.emptySession.AddSampledValue(samples)
	s.Assert().EqualValues(expected, s.emptySession.Consumption)
}

func (s *sessionTestSuite) TestStartSession() {
	currentTime := time.Now()
	var expectedSession = Session{
		IsActive:      true,
		TransactionId: "test1234",
		TagId:         "test1234",
		Started:       &currentTime,
		Consumption:   []types.MeterValue{},
	}

	err := s.emptySession.StartSession("", "")
	s.Assert().Error(err)
	s.Assert().EqualValues(s.emptySession, s.emptySession)

	err = s.emptySession.StartSession("1234", "")
	s.Assert().Error(err)
	s.Assert().EqualValues(s.emptySession, s.emptySession)

	err = s.emptySession.StartSession("", "test1234")
	s.Assert().Error(err)
	s.Assert().EqualValues(s.emptySession, s.emptySession)

	// Ok case
	err = s.emptySession.StartSession("test1234", "test1234")
	s.Assert().NoError(err)
	s.Assert().EqualValues(expectedSession, s.emptySession)
}

func (s *sessionTestSuite) TestEndSession() {
	s.validSession.EndSession()
	s.Assert().EqualValues(s.validSession, s.emptySession)
}

func (s *sessionTestSuite) TestCalculateAvgPower() {
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
	s.Assert().InDelta(0, s.emptySession.CalculateAvgPower(), 0.1)

	s.emptySession.Consumption = consumption1
	s.Assert().InDelta(10, s.emptySession.CalculateAvgPower(), 1)

	s.emptySession.Consumption = mixedConsumption
	s.Assert().InDelta(20, s.emptySession.CalculateAvgPower(), .5)

	s.emptySession.Consumption = faultyMixedConsumption
	s.Assert().InDelta(10, s.emptySession.CalculateAvgPower(), 0.1)

	s.emptySession.Consumption = threeMeasurands
	s.Assert().InDelta(10, s.emptySession.CalculateAvgPower(), 0.1)

}

func (s *sessionTestSuite) TestCalculateEnergyConsumptionWithAvgPower() {
	started5min := time.Now().Add(-5 * time.Minute)

	// Start the session
	err := s.emptySession.StartSession("1234", "1234")
	s.Assert().NoError(err)

	// Add a sample
	s.emptySession.AddSampledValue([]types.SampledValue{
		{
			Value:     "10",
			Measurand: types.MeasurandPowerActiveImport,
		},
	})

	// Session started 5 minutes before
	s.emptySession.Started = &started5min
	s.Assert().InDelta(float64(300*10), s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 8)

	// Session just started
	currentTime := time.Now()
	s.emptySession.Started = &currentTime
	s.Assert().InDelta(0.0, s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 1)

	// Session ended just now
	s.emptySession.EndSession()
	s.Assert().InDelta(0.0, s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 1)

	s.emptySession.AddSampledValue([]types.SampledValue{
		{
			Value:     "10",
			Measurand: types.MeasurandPowerActiveImport,
		},
	})
	s.Assert().InDelta(0.0, s.emptySession.CalculateEnergyConsumptionWithAvgPower(), 1)
}

func TestSession(t *testing.T) {
	suite.Run(t, new(sessionTestSuite))
}
