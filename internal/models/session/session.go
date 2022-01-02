package session

import (
	"errors"
	strUtil "github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
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
		CalculateAvgPower() float32
		CalculateEnergyConsumptionWithAvgPower() float32
		CalculateEnergyConsumption() float32
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

	if !strUtil.IsNotEmpty(tagId) {
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

// AddSampledValue Add all the samples taken to the Session.
func (session *Session) AddSampledValue(samples []types.SampledValue) {
	if session.IsActive {
		log.Tracef("Added meter sample for session %s", session.TransactionId)
		session.Consumption = append(session.Consumption, types.MeterValue{SampledValue: samples})
	}
}

// CalculateAvgPower calculate the average power for a session based on sampled values
func (session *Session) CalculateAvgPower() float32 {
	var (
		powerSum   float32 = 0.0
		numSamples         = 0
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
			case types.MeasurandEnergyActiveExportInterval:
				break
			case types.MeasurandCurrentExport:
				hasCurrent = true
				currentVal, err := strconv.ParseFloat(sampledValue.Value, 32)
				if err != nil {
					continue
				}

				current = currentVal
				break
			case types.MeasurandPowerActiveExport:
				power, err := strconv.ParseFloat(sampledValue.Value, 32)
				if err != nil {
					continue
				}

				hasPower = true
				isValidSample = true
				powerSum += float32(power)
				break
			case types.MeasurandVoltage:
				hasVoltage = true
				voltageVal, err := strconv.ParseFloat(sampledValue.Value, 32)
				if err != nil {
					continue
				}

				voltage = voltageVal
				break
			}
		}
		// If both the current and voltage were sampled and power wasn't, calculate the power by multiplying voltage and current
		if hasCurrent && hasVoltage && !hasPower {
			isValidSample = true
			powerSum += float32(voltage * current)
		}
		// Edge case -> number of samples != length of measurements
		// If there is an array of samples that does not contain both Voltage and Current pair or Power sample, discard the sample
		if isValidSample {
			numSamples++
		}
	}

	if len(session.Consumption) > 0 && numSamples > 0 {
		return powerSum / float32(numSamples)
	}

	return 0
}

// CalculateEnergyConsumptionWithAvgPower calculate the total energy consumption for a session that was active, if it had any measurements
func (session *Session) CalculateEnergyConsumptionWithAvgPower() float32 {
	startDate, err := time.Parse(time.RFC3339, session.Started)
	if err != nil {
		return 0
	}

	var duration = time.Now().Sub(startDate).Seconds()
	// For testing purposes discard any sub 1-second durations
	if duration < 1 {
		return 0
	}

	return session.CalculateAvgPower() * float32(duration)
}

// CalculateEnergyConsumption calculate the total energy consumption for a session that was active only with energy measurments
func (session *Session) CalculateEnergyConsumption() float32 {
	var (
		energySum float32 = 0.0
	)

	for _, meterValue := range session.Consumption {
		for _, sampledValue := range meterValue.SampledValue {
			switch sampledValue.Measurand {
			case types.MeasurandEnergyActiveExportInterval:
				energySample, err := strconv.ParseFloat(sampledValue.Value, 32)
				if err != nil {
					continue
				}

				if energySample > 0 {
					energySum += float32(energySample)
				}

				break
			}
		}
	}

	return energySum
}
