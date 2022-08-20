package session

import (
	"errors"
	strUtil "github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	settingsModel "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"strconv"
	"time"
)

var (
	ErrSessionActive        = errors.New("session already active on the connector")
	ErrInvalidTagId         = errors.New("tag ID invalid")
	ErrInvalidTransactionId = errors.New("transaction id invalid")
)

type (
	Session struct {
		IsActive      bool
		TransactionId string
		TagId         string
		Started       string
		Consumption   []types.MeterValue
	}

	SessionInterface interface {
		StartSession(transactionId string, tagId string) error
		EndSession()
		AddSampledValue(samples []types.SampledValue)
		CalculateAvgPower() float64
		CalculateEnergyConsumptionWithAvgPower() float64
		CalculateEnergyConsumption() float64
		UpdateSessionFile(evseId, connectorId int)
		ToSessionFile() *settingsModel.Session
	}
)

func NewEmptySession() *Session {
	return &Session{
		TransactionId: "",
		TagId:         "",
		IsActive:      false,
		Consumption:   []types.MeterValue{},
	}
}

// StartSession Starts the Session, storing the transactionId and tagId of the user.
// Checks if transaction and tagId are valid strings.
func (session *Session) StartSession(transactionId string, tagId string) error {
	if session.IsActive {
		return ErrSessionActive
	}

	if !strUtil.IsAlphanumeric(transactionId) {
		return ErrInvalidTransactionId
	}

	if strUtil.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	log.Debugf("Started a session %s for %s", transactionId, tagId)

	session.TransactionId = transactionId
	session.TagId = tagId
	session.IsActive = true
	session.Started = time.Now().Format(time.RFC3339)
	session.Consumption = []types.MeterValue{}
	return nil
}

// EndSession End the Session if one is active. Reset the attributes, except the measurands.
func (session *Session) EndSession() {
	if session.IsActive {
		log.Debugf("Ended a session %s for %s", session.TransactionId, session.TagId)
		session.TransactionId = ""
		session.TagId = ""
		session.IsActive = false
		session.Started = ""
	}
}

// ToSessionFile transforms session to a session that is stored in a file in case of a failure.
func (session *Session) ToSessionFile() *settingsModel.Session {
	return &settingsModel.Session{
		IsActive:      session.IsActive,
		TagId:         session.TagId,
		TransactionId: session.TransactionId,
		Started:       session.Started,
		Consumption:   session.Consumption,
	}
}

// UpdateSessionFile Updates the session file.
func (session *Session) UpdateSessionFile(evseId int) {
	settings.UpdateSession(
		evseId,
		session.ToSessionFile(),
	)
}

// AddSampledValue Add all the samples taken to the Session.
func (session *Session) AddSampledValue(samples []types.SampledValue) {
	if session.IsActive {
		log.Tracef("Added meter sample for session %s", session.TransactionId)
		session.Consumption = append(session.Consumption, types.MeterValue{SampledValue: samples})
	}
}

// CalculateAvgPower calculate the average power for a session based on sampled values
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
			sampleValue, err := strconv.ParseFloat(sampledValue.Value, 32)
			if err != nil {
				continue
			}

			switch sampledValue.Measurand {
			case types.MeasurandCurrentImport:
				hasCurrent = true
				current = sampleValue
			case types.MeasurandCurrentExport:
				hasCurrent = true
				current = -sampleValue
			case types.MeasurandPowerActiveImport:
				hasPower = true
				isValidSample = true

				switch sampledValue.Unit {
				case types.UnitOfMeasureKW:
					powerSum += sampleValue * 1000
				case types.UnitOfMeasureW:
					powerSum += sampleValue
				default:
					powerSum += sampleValue
				}

			case types.MeasurandPowerActiveExport:
				hasPower = true
				isValidSample = true

				switch sampledValue.Unit {
				case types.UnitOfMeasureKW:
					powerSum -= sampleValue * 1000
				case types.UnitOfMeasureW:
					powerSum -= sampleValue
				default:
					powerSum -= sampleValue
				}

			case types.MeasurandVoltage:
				hasVoltage = true
				voltage = sampleValue
			}
		}

		// If both the current and voltage were sampled and power wasn't, calculate the power by multiplying voltage and current
		if hasCurrent && hasVoltage && !hasPower {
			isValidSample = true
			powerSum += voltage * current
		}
		// Edge case -> number of samples != length of measurements
		// If there is an array of samples that does not contain both Voltage and Current pair or Power sample, discard the sample
		if isValidSample {
			numSamples++
		}
	}

	if len(session.Consumption) > 0 && numSamples > 0 {
		return powerSum / float64(numSamples)
	}

	return 0
}

// CalculateEnergyConsumptionWithAvgPower calculate the total energy consumption for a session that was active, if it had any measurements
func (session *Session) CalculateEnergyConsumptionWithAvgPower() float64 {
	startDate, err := time.Parse(time.RFC3339, session.Started)
	if err != nil {
		return 0
	}

	var duration = time.Now().Sub(startDate).Seconds()
	// For testing purposes discard any sub 1-second durations
	if duration < 1 {
		return 0
	}

	return session.CalculateAvgPower() * duration
}

// CalculateEnergyConsumption calculate the total energy consumption for a session that was active only with energy measurments
func (session *Session) CalculateEnergyConsumption() float64 {
	var (
		energySum = 0.0
	)

	for _, meterValue := range session.Consumption {
		for _, sampledValue := range meterValue.SampledValue {
			energySample, err := strconv.ParseFloat(sampledValue.Value, 32)
			if err != nil {
				continue
			}

			switch sampledValue.Measurand {
			case types.MeasurandEnergyActiveImportInterval, types.MeasurandEnergyActiveImportRegister:

				switch sampledValue.Unit {
				case types.UnitOfMeasureKWh:
					energySum += energySample * 1000
				case types.UnitOfMeasureWh:
					energySum += energySample
				default:
					energySum += energySample
				}
			}
		}
	}

	return energySum
}
