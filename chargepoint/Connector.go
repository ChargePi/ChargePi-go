package chargepoint

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	power_meter "github.com/xBlaz3kx/ChargePi-go/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"time"
)

type Connector struct {
	EvseId                       int
	ConnectorId                  int
	ConnectorType                string
	ConnectorStatus              core.ChargePointStatus
	ErrorCode                    core.ChargePointErrorCode
	relay                        *hardware.Relay
	powerMeter                   power_meter.PowerMeter
	PowerMeterEnabled            bool
	MaxChargingTime              int
	reservationId                int
	session                      *data.Session
	connectorNotificationChannel chan rxgo.Item
}

// NewConnector Create a new connector object from the provided arguments. EvseId, connectorId and maxChargingTime must be greater than zero.
// When created, it makes an empty session, turns off the relay and defaults the status to Available.
func NewConnector(
	EvseId int,
	connectorId int,
	connectorType string,
	relay *hardware.Relay,
	powerMeter power_meter.PowerMeter,
	powerMeterEnabled bool,
	maxChargingTime int,
) (*Connector, error) {
	if maxChargingTime <= 0 {
		maxChargingTime = 180
	}

	if EvseId <= 0 {
		return nil, errors.New("invalid evse id")
	}

	if connectorId <= 0 {
		return nil, errors.New("invalid connector id")
	}

	if relay == nil {
		return nil, fmt.Errorf("relay pointer cannot be nil")
	}

	relay.Off()
	return &Connector{
		EvseId:            EvseId,
		ConnectorId:       connectorId,
		ConnectorType:     connectorType,
		relay:             relay,
		powerMeter:        powerMeter,
		reservationId:     -1,
		PowerMeterEnabled: powerMeterEnabled,
		MaxChargingTime:   maxChargingTime,
		ConnectorStatus:   core.ChargePointStatusAvailable,
		session: &data.Session{
			TransactionId: "",
			TagId:         "",
			Started:       "",
			Consumption:   []types.MeterValue{},
		},
	}, nil
}

// StartCharging Start charging a connector if connector is available and session could be started.
// It turns on the relay (even if negative logic applies).
func (connector *Connector) StartCharging(transactionId string, tagId string) error {
	var hasSessionStarted = false
	if (connector.IsAvailable() || connector.IsPreparing()) && !connector.session.IsActive {
		hasSessionStarted = connector.session.StartSession(transactionId, tagId)
		connector.SetStatus(core.ChargePointStatusPreparing, core.NoError)
		if hasSessionStarted {
			connector.relay.On()
			connector.SetStatus(core.ChargePointStatusCharging, core.NoError)
			settings.UpdateConnectorSessionInfo(
				connector.EvseId,
				connector.ConnectorId,
				&settings.Session{
					IsActive:      connector.session.IsActive,
					TagId:         connector.session.TagId,
					TransactionId: connector.session.TransactionId,
					Started:       connector.session.Started,
					Consumption:   connector.session.Consumption,
				})
			return nil
		}
		return errors.New("transaction or tag ID invalid")
	}
	return errors.New("invalid connector status or session already active")
}

// ResumeCharging Resumes or restores the charging state after boot if a charging session was active.
func (connector *Connector) ResumeCharging(session data.Session) (err error, chargingTimeElapsed int) {
	//set the transaction id so connector is able to stop the transaction if charging fails
	connector.session.TransactionId = session.TransactionId

	startedChargingTime, err := time.Parse(time.RFC3339, session.Started)
	if err != nil {
		return
	}

	chargingTimeElapsed = int(time.Now().Sub(startedChargingTime).Minutes())
	if connector.MaxChargingTime < chargingTimeElapsed {
		chargingTimeElapsed = connector.MaxChargingTime
		err = fmt.Errorf("session time limit exceeded")
		return
	}

	if connector.IsCharging() || connector.IsPreparing() {
		hasSessionStarted := connector.session.StartSession(session.TransactionId, session.TagId)
		if hasSessionStarted {
			connector.relay.On()
			connector.session.Started = session.Started
			connector.session.Consumption = append(connector.session.Consumption, session.Consumption...)
			return nil, chargingTimeElapsed
		}
		err = errors.New("cannot resume session: unable to start session")
		return
	}

	return errors.New("cannot resume session: invalid connector status"), connector.MaxChargingTime
}

// StopCharging Stops charging the connector by turning the relay off and ending the session.
func (connector *Connector) StopCharging(reason core.Reason) error {
	if connector.IsCharging() || connector.IsPreparing() {
		connector.session.EndSession()
		connector.relay.Off()
		settings.UpdateConnectorSessionInfo(
			connector.EvseId,
			connector.ConnectorId,
			&settings.Session{
				IsActive:      connector.session.IsActive,
				TagId:         connector.session.TagId,
				TransactionId: connector.session.TransactionId,
				Started:       connector.session.Started,
				Consumption:   connector.session.Consumption,
			})
		switch reason {
		case core.ReasonEVDisconnected:
			connector.SetStatus(core.ChargePointStatusSuspendedEVSE, core.NoError)
			break
		case core.ReasonUnlockCommand:
			connector.SetStatus(core.ChargePointStatusUnavailable, core.NoError)
			break
		default:
			connector.SetStatus(core.ChargePointStatusFinishing, core.NoError)
			connector.SetStatus(core.ChargePointStatusAvailable, core.NoError)
		}
		return nil
	}
	return errors.New("connector not charging")
}

// SamplePowerMeter Get a sample from the power meter. The measurands argument takes the list of all the types of the measurands to sample.
// It will add all the samples to the connector's Session if it is active.
func (connector *Connector) SamplePowerMeter(measurands []types.Measurand) {
	if !connector.PowerMeterEnabled || connector.powerMeter == nil {
		return
	}
	log.Println("Sampling connector", connector.ConnectorId)
	var samples []types.SampledValue
	for _, measurand := range measurands {
		var value = 0.0
		switch measurand {
		case types.MeasurandEnergyActiveExportInterval:
			value = connector.powerMeter.GetEnergy()
			break
		case types.MeasurandCurrentExport:
			value = connector.powerMeter.GetCurrent()
			break
		case types.MeasurandPowerActiveExport:
			value = connector.powerMeter.GetPower()
			break
		case types.MeasurandVoltage:
			value = connector.powerMeter.GetVoltage()
			break
		}
		if value != 0.0 {
			samples = append(samples, types.SampledValue{Value: fmt.Sprintf("%.3f", value), Measurand: measurand})
		}
	}
	connector.session.AddSampledValue(samples)
}

// preparePowerMeterAtConnector
func (connector *Connector) preparePowerMeterAtConnector() error {
	var (
		measurands []types.Measurand
		err        error
	)
	cache.Cache.Set(fmt.Sprintf("MeterValueLastIndex%d%d", connector.EvseId, connector.ConnectorId),
		0, time.Duration(connector.MaxChargingTime)*time.Minute)
	measurands = getTypesToSample()
	// Get the sample interval
	sampleInterval, err := settings.GetConfigurationValue("MeterValueSampleInterval")
	if err != nil {
		sampleInterval = "10"
	}
	// schedule the sampling
	_, err = scheduler.GetScheduler().Every(fmt.Sprintf("%ss", sampleInterval)).
		Tag(fmt.Sprintf("connector%dSampling", connector.ConnectorId)).Do(connector.SamplePowerMeter, measurands)
	if err != nil {
		return err
	}
	return nil
}

func (connector *Connector) IsAvailable() bool {
	return connector.ConnectorStatus == core.ChargePointStatusAvailable
}
func (connector *Connector) IsCharging() bool {
	return connector.ConnectorStatus == core.ChargePointStatusCharging
}
func (connector *Connector) IsPreparing() bool {
	return connector.ConnectorStatus == core.ChargePointStatusPreparing
}
func (connector *Connector) IsReserved() bool {
	return connector.ConnectorStatus == core.ChargePointStatusReserved
}
func (connector *Connector) IsUnavailable() bool {
	return connector.ConnectorStatus == core.ChargePointStatusUnavailable
}

func (connector *Connector) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	time.Sleep(time.Millisecond * 100)
	connector.ConnectorStatus = status
	connector.ErrorCode = errCode
	settings.UpdateConnectorStatus(connector.EvseId, connector.ConnectorId, status)
	time.Sleep(time.Millisecond * 100)
	if connector.connectorNotificationChannel != nil {
		connector.connectorNotificationChannel <- rxgo.Of(connector)
	}
}

func (connector *Connector) GetTransactionId() string {
	return connector.session.TransactionId
}
func (connector *Connector) GetTagId() string {
	return connector.session.TagId
}

func (connector *Connector) ReserveConnector(reservationId int) error {
	if reservationId <= 0 {
		return fmt.Errorf("reservation id is invalid")
	}
	if !connector.IsAvailable() {
		return fmt.Errorf("connector status invalid")
	}
	connector.reservationId = reservationId
	connector.SetStatus(core.ChargePointStatusReserved, core.NoError)
	return nil
}
func (connector *Connector) RemoveReservation() error {
	if !connector.IsReserved() {
		return fmt.Errorf("connector status invalid")
	}
	connector.reservationId = -1
	connector.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	return nil
}

func (connector *Connector) GetReservationId() int {
	return connector.reservationId
}
