package util

import (
	"math"
	"reflect"
	"strconv"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"
)

// IsNilInterfaceOrPointer check if the variable is nil or if the pointer's value is nil.
func IsNilInterfaceOrPointer(sth interface{}) bool {
	return sth == nil || (reflect.ValueOf(sth).Kind() == reflect.Ptr && reflect.ValueOf(sth).IsNil())
}

func ValidateMeterValueSample(sample types.SampledValue) error {
	sampleValue, err := strconv.ParseFloat(sample.Value, 32)
	if err != nil {
		return errors.Wrap(err, "error parsing sample value")
	}

	// todo sanity check for other measurands
	switch sample.Measurand {
	case types.MeasurandEnergyActiveImportRegister,
		types.MeasurandEnergyActiveExportInterval,
		types.MeasurandEnergyActiveImportInterval,
		types.MeasurandEnergyActiveExportRegister:
		// Must not be negative
		if sampleValue < 0 {
			return errors.New("energy cannot be negative")
		}

		// If unit is Wh, must be a full number
		if sample.Unit == types.UnitOfMeasureWh || sample.Unit == "" && sampleValue != math.Trunc(sampleValue) {
			return errors.New("energy must be a full number")
		}

	case types.MeasurandPowerActiveImport:
		// Cannot be negative, as that would mean exporting power
		if sampleValue < 0 {
			return errors.New("active power import cannot be negative")
		}

	case types.MeasurandPowerActiveExport:
		// Cannot be positive, as that would mean importing power?
		if sampleValue > 0 {
			return errors.New("active power export cannot be positive")
		}
	case types.MeasurandFrequency:
		// Must be positive
		if sampleValue <= 0 {
			return errors.New("frequency must be positive")
		}
	case types.MeasurandRPM:
		// Must be positive & full number
		if sampleValue < 0 {
			return errors.New("RPM must be positive")
		}

		if sampleValue != math.Trunc(sampleValue) {
			return errors.New("RPM must be a full number")
		}
	default:
		return errors.New("unknown measurand")
	}

	return nil
}
