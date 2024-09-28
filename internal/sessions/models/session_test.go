package models

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
)

type sessionTestSuite struct {
	suite.Suite
}

func (s *sessionTestSuite) SetupTest() {

}

func (s *sessionTestSuite) TestAddSampledValue() {
	tests := []struct {
		name            string
		samples         []types.SampledValue
		expectedSamples []types.SampledValue
		isSessionActive bool
		err             error
		anyErr          bool
	}{
		{
			name: "Add sample to inactive session",
			samples: []types.SampledValue{
				{
					Value:     "123.21",
					Measurand: types.MeasurandCurrentImport,
				},
				{
					Value:     "123.21",
					Measurand: types.MeasurandVoltage,
				},
				{
					Value:     "123",
					Measurand: types.MeasurandPowerActiveImport,
				},
			},
			expectedSamples: []types.SampledValue{},
			isSessionActive: false,
			err:             ErrSessionNotActive,
		},
		{
			name: "Add sample to active session",
			samples: []types.SampledValue{
				{
					Value:     "123.21",
					Measurand: types.MeasurandCurrentImport,
				},
				{
					Value:     "123.21",
					Measurand: types.MeasurandVoltage,
				},
				{
					Value:     "123",
					Measurand: types.MeasurandPowerActiveImport,
				},
			},
			expectedSamples: []types.SampledValue{
				{
					Value:     "123.21",
					Measurand: types.MeasurandCurrentImport,
				},
				{
					Value:     "123.21",
					Measurand: types.MeasurandVoltage,
				},
				{
					Value:     "123",
					Measurand: types.MeasurandPowerActiveImport,
				}},
			isSessionActive: true,
			err:             nil,
		},
		{
			name: "Wrong measurand",
			samples: []types.SampledValue{
				{
					Value:     "123.21",
					Measurand: "unknown",
				},
			},
			isSessionActive: true,
			expectedSamples: []types.SampledValue{},
			err:             nil,
			anyErr:          true,
		},
		{
			name: "Add sample to active session with wrong value",
			samples: []types.SampledValue{
				{
					Value:     "-123.21",
					Measurand: types.MeasurandPowerActiveImport,
				},
			},
			isSessionActive: true,
			expectedSamples: []types.SampledValue{},
			err:             nil,
			anyErr:          true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			session := NewEmptySession()
			session.IsActive = tt.isSessionActive

			err := session.AddSampledValue(tt.samples)
			if tt.err != nil {
				s.ErrorIs(err, tt.err)
			} else if tt.anyErr {
				s.Error(err)
			} else {
				s.NoError(err)
				// Check if the samples were added
				s.Len(session.Consumption, 1)
				s.ElementsMatch(tt.expectedSamples, session.Consumption[0].SampledValue)
			}
		})
	}
}

func (s *sessionTestSuite) TestStartSession() {
	tests := []struct {
		name  string
		tagId string
		err   error
	}{
		{
			name:  "Start session",
			tagId: "test1234",
			err:   nil,
		},
		{
			name:  "Start already active session",
			tagId: "test1234",
			err:   ErrSessionActive,
		},
		{
			name:  "Start session with invalid tag ID",
			tagId: "",
			err:   ErrInvalidTagId,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			session := NewEmptySession()

			if tt.name == "Start already active session" {
				session.IsActive = true
			}

			err := session.StartSession(tt.tagId)
			if tt.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
				s.True(session.IsActive)
				s.Equal(tt.tagId, session.TagId)
				s.NotNil(session.Started)
			}
		})
	}
}

func (s *sessionTestSuite) TestEndSession() {
	tests := []struct {
		name    string
		session Session
		err     error
	}{
		{
			name: "End inactive session",
			session: Session{
				IsActive: false,
			},
			err: ErrSessionNotActive,
		},
		{
			name: "End active session",
			session: Session{
				IsActive: true,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := tt.session.EndSession()
			if tt.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
				s.False(tt.session.IsActive)
			}
		})
	}
}

func (s *sessionTestSuite) TestSetTransactionId() {
	tests := []struct {
		name          string
		transactionId string
		isActive      bool
		err           error
	}{
		{
			name:          "Invalid transaction ID",
			transactionId: "",
			isActive:      true,
			err:           ErrInvalidTransactionId,
		},
		{
			name:          "Valid transaction ID",
			transactionId: "test1234",
			isActive:      true,
			err:           nil,
		},
		{
			name:          "Session not active",
			transactionId: "test1234",
			isActive:      false,
			err:           ErrSessionNotActive,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			session := NewEmptySession()
			session.IsActive = tt.isActive

			err := session.SetTransactionId(tt.transactionId)
			if tt.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.transactionId, session.TransactionId)
			}
		})
	}
}

func (s *sessionTestSuite) TestGetMeterValues() {
	tests := []struct {
		name     string
		samples  []types.SampledValue
		expected []types.MeterValue
	}{
		{
			name: "Single sample",
			samples: []types.SampledValue{
				{
					Value:     "10",
					Unit:      types.UnitOfMeasureKW,
					Measurand: types.MeasurandPowerActiveImport,
				},
			},
			expected: []types.MeterValue{
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "10",
							Unit:      types.UnitOfMeasureKW,
							Measurand: types.MeasurandPowerActiveImport,
						},
					},
				},
			},
		},
		{
			name: "Multiple samples",
			samples: []types.SampledValue{
				{
					Value:     "10",
					Unit:      types.UnitOfMeasureKW,
					Measurand: types.MeasurandPowerActiveImport,
				},
			},
			expected: []types.MeterValue{
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "10",
							Unit:      types.UnitOfMeasureKW,
							Measurand: types.MeasurandPowerActiveImport,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			session := NewEmptySession()
			err := session.StartSession("test1234")
			s.Require().NoError(err)

			err = session.AddSampledValue(tt.samples)
			s.Require().NoError(err)

			meterValues := session.GetMeterValues()
			for i, value := range meterValues {
				s.InDelta(tt.expected[i].Timestamp.Second(), value.Timestamp.Second(), 1.0)
				s.ElementsMatch(tt.expected[i].SampledValue, value.SampledValue)
			}
		})
	}
}

func (s *sessionTestSuite) TestCalculateAvgPower() {
	emptySample := []types.SampledValue{}

	tests := []struct {
		name        string
		meterValues []types.MeterValue
		expected    float64
	}{
		{
			name: "Empty samples",
			meterValues: []types.MeterValue{
				{
					Timestamp:    types.Now(),
					SampledValue: emptySample,
				}, {
					Timestamp:    types.Now(),
					SampledValue: emptySample,
				},
			},
			expected: 0,
		},
		{
			name: "Single meter value",
			meterValues: []types.MeterValue{
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "10",
							Unit:      types.UnitOfMeasureKW,
							Measurand: types.MeasurandPowerActiveImport,
						},
					},
				},
			},
			expected: 10000,
		},
		{
			name: "Mixed samples",
			meterValues: []types.MeterValue{
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "10",
							Unit:      types.UnitOfMeasureKW,
							Measurand: types.MeasurandPowerActiveImport,
						},
					},
				},
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "5",
							Unit:      types.UnitOfMeasureKW,
							Measurand: types.MeasurandPowerActiveImport,
						},
					},
				},
			},
			expected: 7500.0,
		},
		{
			name: "Faulty mixed samples",
			meterValues: []types.MeterValue{
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "10",
							Unit:      types.UnitOfMeasureKW,
							Measurand: types.MeasurandPowerActiveImport,
						},
					},
				},
				{
					Timestamp: types.Now(),
					SampledValue: []types.SampledValue{
						{
							Value:     "2",
							Measurand: types.MeasurandCurrentImport,
						},
					},
				},
			},
			expected: 10000,
		},
		{
			name: "Two samples with three measurands",
			meterValues: []types.MeterValue{
				{
					Timestamp: types.Now(),
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
				},
				{
					Timestamp: types.Now(),
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
			},
			expected: 10,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			session := NewEmptySession()
			session.Consumption = tt.meterValues

			avgPower := session.CalculateAvgPower()
			s.Assert().InDelta(tt.expected, avgPower, 0.2)
		})
	}
}

func (s *sessionTestSuite) TestGetEnergyConsumption() {
	started5min := time.Now().Add(-5 * time.Minute)

	tests := []struct {
		name           string
		session        *Session
		expectedEnergy float64
		expectErr      bool
	}{
		{
			name:           "Empty session",
			session:        NewEmptySession(),
			expectedEnergy: 0.0,
		},
		{
			name: "Duration less than a second",
			session: &Session{
				Started:     lo.ToPtr(time.Now().Add(-100 * time.Millisecond)),
				Consumption: []types.MeterValue{},
			},
			expectedEnergy: 0.0,
		},
		{
			name: "Single meter value",
			session: &Session{
				Started: lo.ToPtr(time.Now().Add(-2 * time.Minute)),
				Consumption: []types.MeterValue{
					{
						Timestamp: types.NewDateTime(time.Now().Add(-1 * time.Minute)),
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Unit:      types.UnitOfMeasureKWh,
								Measurand: types.MeasurandEnergyActiveImportInterval,
							},
						},
					},
				},
			},
			expectedEnergy: 1000.0,
		},
		{
			name: "Multiple meter values",
			session: &Session{
				Started: &started5min,
				Consumption: []types.MeterValue{
					{
						Timestamp: types.NewDateTime(started5min),
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Unit:      types.UnitOfMeasureKWh,
								Measurand: types.MeasurandEnergyActiveImportInterval,
							},
						},
					},
					{
						Timestamp: types.NewDateTime(time.Now()),
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Unit:      types.UnitOfMeasureKWh,
								Measurand: types.MeasurandEnergyActiveImportInterval,
							},
						},
					},
				},
			},
			expectedEnergy: 2000.0,
		},
		{
			name: "Multiple meter values with different measurands",
			session: &Session{
				Started: &started5min,
				Consumption: []types.MeterValue{
					{
						Timestamp: types.NewDateTime(started5min),
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Unit:      types.UnitOfMeasureKWh,
								Measurand: types.MeasurandEnergyActiveImportInterval,
							},
						},
					},
					{
						Timestamp: types.NewDateTime(time.Now()),
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Unit:      types.UnitOfMeasureKWh,
								Measurand: types.MeasurandEnergyActiveImportInterval,
							},
						},
					},
					{
						Timestamp: types.NewDateTime(time.Now()),
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Unit:      types.UnitOfMeasureKW,
								Measurand: types.MeasurandPowerActiveImport,
							},
						},
					},
				},
			},
			expectedEnergy: 2000.0,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			energy, err := tt.session.GetEnergyConsumption()
			if tt.expectErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.InDelta(tt.expectedEnergy, energy, 0.2)
		})
	}
}

func (s *sessionTestSuite) TestGetMeterValuesFrom() {
	referenceTimeNow := time.Now()

	tests := []struct {
		name        string
		from        time.Time
		meterValues []types.MeterValue
		expected    []types.MeterValue
	}{
		{
			name:        "Empty session",
			from:        time.Now().Add(-1 * time.Hour),
			meterValues: []types.MeterValue{},
			expected:    []types.MeterValue{},
		},
		{
			name: "All samples in date range",
			meterValues: []types.MeterValue{
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-1 * time.Hour)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-30 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-10 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-5 * time.Minute)),
				},
			},
			expected: []types.MeterValue{
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-1 * time.Hour)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-30 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-10 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-5 * time.Minute)),
				},
			},
		},
		{
			name: "Some samples in date range",
			from: time.Now().Add(-20 * time.Minute),
			meterValues: []types.MeterValue{
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-1 * time.Hour)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-30 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-10 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-5 * time.Minute)),
				},
			},
			expected: []types.MeterValue{
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-10 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-5 * time.Minute)),
				},
			},
		},
		{
			name: "Meter value not in date range",
			from: time.Now().Add(-1 * time.Minute),
			meterValues: []types.MeterValue{
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-1 * time.Hour)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-30 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-10 * time.Minute)),
				},
				{
					Timestamp: types.NewDateTime(referenceTimeNow.Add(-5 * time.Minute)),
				},
			},
			expected: []types.MeterValue{},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			session := NewEmptySession()
			session.Consumption = tt.meterValues

			meterValues := session.GetMeterValuesFrom(tt.from)
			s.Require().Len(meterValues, len(tt.expected))

			for i, value := range tt.expected {
				s.EqualValues(*value.Timestamp, *meterValues[i].Timestamp)
			}
		})
	}
}

func TestSession(t *testing.T) {
	suite.Run(t, new(sessionTestSuite))
}
