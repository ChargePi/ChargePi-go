package util

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
)

// nolint:all
type exampleInterface interface {
	ExampleMethod()
}

type exampleStructImplementsInterface struct {
}

func (e exampleStructImplementsInterface) ExampleMethod() {
}

func TestIsNilInterfaceOrPointer(t *testing.T) {
	tests := []struct {
		name  string
		thing interface{}
		want  bool
	}{
		{
			name:  "nil interface",
			thing: exampleInterface(nil),
			want:  true,
		},
		{
			name:  "nil pointer",
			thing: nil,
			want:  true,
		},
		{
			name:  "valid interface",
			thing: &exampleStructImplementsInterface{},
			want:  false,
		},
		{
			name:  "valid pointer",
			thing: new(int),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsNilInterfaceOrPointer(tt.thing))
		})
	}
}

func TestValidateMeterValueSample(t *testing.T) {
	tests := []struct {
		name    string
		sample  types.SampledValue
		wantErr bool
		error   string
	}{
		{
			name: "Wrong measurand",
			sample: types.SampledValue{
				Value:     "123.21",
				Measurand: "unknown",
			},
			wantErr: true,
			error:   "unknown measurand",
		},
		{
			name: "Invalid value for power",
			sample: types.SampledValue{
				Value:     "-123.21",
				Measurand: types.MeasurandPowerActiveImport,
			},
			wantErr: true,
			error:   "active power import cannot be negative",
		},
		{
			name: "Valid imported power",
			sample: types.SampledValue{
				Value:     "123.21",
				Measurand: types.MeasurandPowerActiveImport,
			},
			wantErr: false,
		},
		{
			name: "Invalid value for energy",
			sample: types.SampledValue{
				Value:     "-123",
				Measurand: types.MeasurandEnergyActiveImportInterval,
			},
			wantErr: true,
			error:   "energy cannot be negative",
		},
		{
			name: "Invalid value for energy - cannot be a fraction",
			sample: types.SampledValue{
				Value:     "123.21",
				Measurand: types.MeasurandEnergyActiveImportInterval,
			},
			wantErr: true,
			error:   "energy must be a full number",
		},
		{
			name: "Invalid value for energy - cannot be a fraction",
			sample: types.SampledValue{
				Value:     "123.21",
				Unit:      types.UnitOfMeasureWh,
				Measurand: types.MeasurandEnergyActiveImportInterval,
			},
			wantErr: true,
			error:   "energy must be a full number",
		},
		{
			name: "Energy allowed to be fractional if unit is Wh",
			sample: types.SampledValue{
				Value:     "123.21",
				Unit:      types.UnitOfMeasureKWh,
				Measurand: types.MeasurandEnergyActiveImportInterval,
			},
			wantErr: false,
		},
		{
			name: "Valid value for energy",
			sample: types.SampledValue{
				Value:     "1234",
				Measurand: types.MeasurandEnergyActiveImportInterval,
			},
			wantErr: false,
		},
		{
			name: "Invalid value for frequency",
			sample: types.SampledValue{
				Value:     "-1",
				Measurand: types.MeasurandFrequency,
			},
			wantErr: true,
			error:   "frequency must be positive",
		},
		{
			name: "Valid value for frequency",
			sample: types.SampledValue{
				Value:     "50",
				Measurand: types.MeasurandFrequency,
			},
			wantErr: false,
		},
		{
			name: "Invalid value for RPM - negative",
			sample: types.SampledValue{
				Value:     "-123",
				Measurand: types.MeasurandRPM,
			},
			wantErr: true,
			error:   "RPM must be positive",
		},
		{
			name: "Invalid value for RPM - fraction",
			sample: types.SampledValue{
				Value:     "123.21",
				Measurand: types.MeasurandRPM,
			},
			wantErr: true,
			error:   "RPM must be a full number",
		},
		{
			name: "Valid value for RPM",
			sample: types.SampledValue{
				Value:     "123",
				Measurand: types.MeasurandRPM,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMeterValueSample(tt.sample)
			if tt.wantErr {
				assert.ErrorContains(t, err, tt.error)
				return
			}

			assert.NoError(t, err)
		})
	}
}
