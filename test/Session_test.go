package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"testing"
	"time"
)

func TestSession_AddSampledValue(t *testing.T) {
	require := require.New(t)

	samples := []types.SampledValue{
		{
			Value:     "123.21",
			Measurand: types.MeasurandCurrentExport,
		}, {
			Value:     "123.21",
			Measurand: types.MeasurandVoltage,
		}, {
			Value:     "123.21",
			Measurand: types.MeasurandPowerActiveExport,
		},
	}

	session := data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	//session not active
	session.AddSampledValue(samples)
	require.Nil(session.Consumption)

	//start session
	hasStarted := session.StartSession("1234", "1234")
	assert.True(t, hasStarted)

	//session active
	session.AddSampledValue(samples)
	expected := []types.MeterValue{
		{
			SampledValue: samples,
		},
	}

	require.EqualValues(expected, session.Consumption)

}

func TestSession_EndSession(t *testing.T) {
	require := require.New(t)

	validSession := data.Session{
		IsActive:      true,
		TransactionId: "1234",
		TagId:         "1234",
		Started:       time.Now().Format(time.RFC3339),
		Consumption:   nil,
	}

	validSession.EndSession()
	require.EqualValues(validSession, data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	})
}

func TestSession_StartSession(t *testing.T) {
	require := require.New(t)

	session := data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	emptySession := data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	session.StartSession("", "")
	require.EqualValues(emptySession, session)

	session.StartSession("1234", "")
	require.EqualValues(emptySession, session)

	session.StartSession("", "test1234")
	require.EqualValues(emptySession, session)

	// ok case
	session.StartSession("test1234", "test1234")
	require.EqualValues(data.Session{
		IsActive:      true,
		TransactionId: "test1234",
		TagId:         "test1234",
		Started:       time.Now().Format(time.RFC3339),
		Consumption:   []types.MeterValue{},
	}, session)
}

func TestSession_CalculateAvgPower(t *testing.T) {
	require := require.New(t)

	zeroSamples := []types.MeterValue{
		{
			Timestamp:    &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{},
		}, {
			Timestamp:    &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{},
		},
	}

	meterWithSamplePower10 := types.MeterValue{
		Timestamp: &types.DateTime{Time: time.Now()},
		SampledValue: []types.SampledValue{
			{
				Value:     "10",
				Measurand: types.MeasurandPowerActiveExport,
			},
		},
	}
	consumption1 := []types.MeterValue{
		meterWithSamplePower10, meterWithSamplePower10,
	}

	mixedConsumption := []types.MeterValue{
		meterWithSamplePower10,
		{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "2",
					Measurand: types.MeasurandCurrentExport,
				}, {
					Value:     "15",
					Measurand: types.MeasurandVoltage,
				},
			},
		},
	}

	faultyMixedConsumption := []types.MeterValue{
		meterWithSamplePower10,
		{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "2",
					Measurand: types.MeasurandCurrentExport,
				},
			},
		},
	}

	threeMeasurands := []types.MeterValue{
		{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "10",
					Measurand: types.MeasurandPowerActiveExport,
				},
				{
					Value:     "2",
					Measurand: types.MeasurandCurrentExport,
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
					Measurand: types.MeasurandCurrentExport,
				},
				{
					Value:     "15",
					Measurand: types.MeasurandVoltage,
				},
				{
					Value:     "10",
					Measurand: types.MeasurandPowerActiveExport,
				},
			},
		},
	}

	session := data.Session{
		IsActive:      true,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   zeroSamples,
	}
	// start with zero samples
	require.InDelta(0, session.CalculateAvgPower(), 0.1)

	session.Consumption = consumption1
	require.InDelta(10, session.CalculateAvgPower(), 1)

	session.Consumption = mixedConsumption
	require.InDelta(20, session.CalculateAvgPower(), .5)

	session.Consumption = faultyMixedConsumption
	require.InDelta(10, session.CalculateAvgPower(), 0.1)

	session.Consumption = threeMeasurands
	require.InDelta(10, session.CalculateAvgPower(), 0.1)
}

func TestSession_CalculateEnergyConsumptionWithAvgPower(t *testing.T) {
	require := require.New(t)

	session := data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}
	started5min := time.Now().Add(-5 * time.Minute).Format(time.RFC3339)

	session.StartSession("1234", "1234")
	session.AddSampledValue([]types.SampledValue{
		{
			Value:     "10",
			Measurand: types.MeasurandPowerActiveExport,
		},
	})

	// session started 5 minutes before
	session.Started = started5min
	require.InDelta(float32(300*10), session.CalculateEnergyConsumptionWithAvgPower(), 8)

	// session just started
	session.Started = time.Now().Format(time.RFC3339)
	require.InDelta(0.0, session.CalculateEnergyConsumptionWithAvgPower(), 1)

	// session ended just now
	session.EndSession()
	require.InDelta(0.0, session.CalculateEnergyConsumptionWithAvgPower(), 1)

	session.AddSampledValue([]types.SampledValue{
		{
			Value:     "10",
			Measurand: types.MeasurandPowerActiveExport,
		},
	})
	require.InDelta(0.0, session.CalculateEnergyConsumptionWithAvgPower(), 1)
}

func TestSession_CalculateEnergyConsumption(t *testing.T) {
	require := require.New(t)

	session := data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	energyConsumption := []types.MeterValue{
		{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "1",
					Measurand: types.MeasurandEnergyActiveExportInterval,
				},
			},
		},
	}

	noEnergyConsumed := []types.MeterValue{
		{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "10",
					Measurand: types.MeasurandVoltage,
				}, {
					Value:     "10",
					Measurand: types.MeasurandCurrentExport,
				}, {
					Value:     "10",
					Measurand: types.MeasurandPowerActiveExport,
				},
			},
		},
	}

	multipleEnergySamples := []types.MeterValue{
		{
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "1",
					Measurand: types.MeasurandEnergyActiveExportInterval,
				},
			},
		}, {
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "10",
					Measurand: types.MeasurandEnergyActiveExportInterval,
				},
			},
		}, {
			Timestamp: &types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{
				{
					Value:     "11",
					Measurand: types.MeasurandEnergyActiveExportInterval,
				},
			},
		},
	}

	session.Consumption = energyConsumption
	require.InDelta(1, session.CalculateEnergyConsumption(), 1)

	session.Consumption = noEnergyConsumed
	require.InDelta(0.0, session.CalculateEnergyConsumption(), 1)

	session.Consumption = multipleEnergySamples
	require.InDelta(22.0, session.CalculateEnergyConsumption(), 1)
}
