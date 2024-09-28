package models

import (
	"time"

	"github.com/ChargePi/ChargePi-go/pkg/util"
	strUtil "github.com/agrison/go-commons-lang/stringUtils"
	"github.com/google/uuid"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

var (
	ErrSessionActive        = errors.New("session already active on the connector")
	ErrSessionNotActive     = errors.New("session is not active")
	ErrInvalidTagId         = errors.New("tag ID invalid")
	ErrInvalidTransactionId = errors.New("transaction id invalid")
	ErrNoSamples            = errors.New("no samples to add")
)

type Session struct {
	// ID is a randomly generated UUID for the session.
	ID uuid.UUID `json:"id"`

	// IsActive indicates if the session is active or not.
	IsActive bool `json:"isActive"`

	// EvseId is the ID of the EVSE the session is associated with.
	EvseId int `json:"evseId"`

	// ConnectorId is the ID of the connector the session is associated with.
	ConnectorId *int `json:"connectorId"`

	// ReservationId is the ID of the reservation the session is associated with.
	ReservationId *int `json:"reservationId"`

	// TransactionId is the OCPP transaction ID associated with the session.
	TransactionId string `json:"transactionId"`

	// TagId is the ID of the tag associated with the session.
	TagId string `json:"tagId"`

	// Started is the time the session was started.
	Started *time.Time `json:"started"`

	// Ended is the time the session was ended.
	Ended *time.Time `json:"ended"`

	// Consumption is the time the session was ended.
	Consumption []types.MeterValue `json:"consumption"`
}

func NewEmptySession() *Session {
	return &Session{
		ID:            uuid.New(),
		TransactionId: "",
		TagId:         "",
		IsActive:      false,
		Consumption:   []types.MeterValue{},
	}
}

// StartSession Starts the Session, storing the transactionId and tagId of the user.
// Checks if transaction and tagId are valid strings.
func (session *Session) StartSession(tagId string) error {
	if strUtil.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	if session.IsActive {
		return ErrSessionActive
	}

	session.ID = uuid.New()
	session.TagId = tagId
	session.IsActive = true
	session.Started = lo.ToPtr(time.Now())
	session.Consumption = []types.MeterValue{}
	return nil
}

func (session *Session) GetTransactionId() string {
	return session.TransactionId
}

func (session *Session) SetTransactionId(transactionId string) error {
	if !strUtil.IsAlphanumeric(transactionId) {
		return ErrInvalidTransactionId
	}

	if !session.IsActive {
		return ErrSessionNotActive
	}

	session.TransactionId = transactionId
	return nil
}

// EndSession End the Session if one is active. Reset the attributes, except the measurands.
func (session *Session) EndSession() error {
	if session.IsActive {
		log.Debugf("Ended a session %s for %s", session.TransactionId, session.TagId)
		session.TransactionId = ""
		session.TagId = ""
		session.IsActive = false
		session.Ended = lo.ToPtr(time.Now())
		return nil
	}

	return ErrSessionNotActive
}

// AddSampledValue Add all the samples taken to the Session.
func (session *Session) AddSampledValue(samples []types.SampledValue) error {
	if len(samples) == 0 {
		return ErrNoSamples
	}

	if !session.IsActive {
		return ErrSessionNotActive
	}

	// Validate samples individually
	for _, sample := range samples {
		err := util.ValidateMeterValueSample(sample)
		if err != nil {
			return errors.Wrap(err, "error validating sample")
		}
	}

	log.Tracef("Added meter sample for session %s", session.TransactionId)
	session.Consumption = append(session.Consumption, types.MeterValue{
		Timestamp:    types.NewDateTime(time.Now()),
		SampledValue: samples,
	})
	return nil
}

// GetMeterValues Get the meter values for session
func (session *Session) GetMeterValues() []types.MeterValue {
	return session.Consumption
}

// GetMeterValuesFrom get the consumption from a specific date
func (session *Session) GetMeterValuesFrom(fromDate time.Time) []types.MeterValue {
	return lo.Filter(session.Consumption, func(meterValue types.MeterValue, _ int) bool {
		return meterValue.Timestamp.Time.After(fromDate)
	})
}

// CalculateAvgPower calculate the average power for a session based on sampled values.
// It will prioritize power samples (if present) and calculate the average power based on them,
// otherwise it will calculate the power based on current and voltage samples.
func (session *Session) CalculateAvgPower() float64 {
	var (
		powerSum   = 0.0
		numSamples = 0
	)

	for _, meterValue := range session.Consumption {
		var (
			hasCurrent    = false
			hasVoltage    = false
			hasPower      = false
			isValidSample = false
			voltage       = 0.0
			current       = 0.0
		)

		for _, sampledValue := range meterValue.SampledValue {
			switch sampledValue.Measurand {
			case types.MeasurandCurrentImport:
				hasCurrent = true
				value, err := util.ToBaseValue(sampledValue)
				if err != nil {
					continue
				}
				current = value

			case types.MeasurandCurrentExport:
				hasCurrent = true
				value, err := util.ToBaseValue(sampledValue)
				if err != nil {
					continue
				}
				current = -value

			case types.MeasurandPowerActiveImport:
				hasPower = true
				isValidSample = true

				value, err := util.ToBaseValue(sampledValue)
				if err != nil {
					continue
				}
				powerSum += value

			case types.MeasurandPowerActiveExport:
				hasPower = true
				isValidSample = true

				value, err := util.ToBaseValue(sampledValue)
				if err != nil {
					continue
				}
				powerSum -= value

			case types.MeasurandVoltage:
				hasVoltage = true
				value, err := util.ToBaseValue(sampledValue)
				if err != nil {
					continue
				}
				voltage = value
			}
		}

		// If both the current and voltage were sampled and power was not,
		// calculate the power by multiplying voltage and current.
		if !hasPower && hasCurrent && hasVoltage {
			isValidSample = true
			powerSum += voltage * current
		}

		// Edge case -> number of samples != length of measurements
		// If there is an array of samples that does not contain both Voltage and Current pair or Power sample,
		// discard the sample
		if isValidSample {
			numSamples++
		}
	}

	if len(session.Consumption) > 0 && numSamples > 0 {
		return powerSum / float64(numSamples)
	}

	return 0
}

// GetEnergyConsumption calculates the total energy consumption for a session that was active.
// The energy consumption is calculated based on measurements - if energy measurements were taken,
// the total energy is calculated by summing the energy measurements or
// substituting with the average power and sample duration (to approximate) if no energy measurements were taken.
// Always returns the value in Wh.
func (session *Session) GetEnergyConsumption() (float64, error) {
	if util.IsNilInterfaceOrPointer(session.Started) {
		return 0, nil
	}

	var sum = 0.0
	var previousMeasurmentTime *types.DateTime

	for i, meterValue := range session.Consumption {
		if i == 0 {
			previousMeasurmentTime = types.NewDateTime(*session.Started)
		}

		// Find the energy measurement and sum it
		energy, wasFound := lo.Find(meterValue.SampledValue, func(sampledValue types.SampledValue) bool {
			switch sampledValue.Measurand {
			case types.MeasurandEnergyActiveImportRegister,
				types.MeasurandEnergyActiveImportInterval,
				types.MeasurandEnergyActiveExportRegister,
				types.MeasurandEnergyActiveExportInterval:
				return true
			default:
				return false
			}
		})
		if wasFound {
			value, err := util.ToBaseValue(energy)
			if err != nil {
				return 0, err
			}

			sum += value
			previousMeasurmentTime = meterValue.Timestamp
			continue
		}

		// If no energy measurement is found, calculate the power and sum it
		power, powerFound := lo.Find(meterValue.SampledValue, func(sampledValue types.SampledValue) bool {
			switch sampledValue.Measurand {
			case types.MeasurandPowerActiveImport,
				types.MeasurandPowerActiveExport:
				return true
			default:
				return false
			}
		})
		if powerFound {
			value, err := util.ToBaseValue(power)
			if err != nil {
				return 0, err
			}

			if previousMeasurmentTime != nil {
				sum += value * meterValue.Timestamp.Sub(previousMeasurmentTime.Time).Seconds()
			}

			previousMeasurmentTime = meterValue.Timestamp
			continue
		}

		// If no power is found, calculate the power based on voltage and current
		voltage, voltageFound := lo.Find(meterValue.SampledValue, func(sampledValue types.SampledValue) bool {
			return sampledValue.Measurand == types.MeasurandVoltage
		})
		current, currentFound := lo.Find(meterValue.SampledValue, func(sampledValue types.SampledValue) bool {
			return sampledValue.Measurand == types.MeasurandCurrentImport
		})
		if voltageFound && currentFound {
			voltageValue, err := util.ToBaseValue(voltage)
			if err != nil {
				return 0, err
			}

			currentValue, err := util.ToBaseValue(current)
			if err != nil {
				return 0, err
			}

			if previousMeasurmentTime != nil {
				sum += voltageValue * currentValue * meterValue.Timestamp.Sub(previousMeasurmentTime.Time).Seconds()
			}

			previousMeasurmentTime = meterValue.Timestamp
		}
	}

	return sum, nil
}
