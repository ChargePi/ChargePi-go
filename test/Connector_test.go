package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"testing"
	"time"
)

func TestConnector_ResumeCharging(t *testing.T) {
	require := require.New(t)

	connector, err := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(11, false), nil, false, 15)
	assert.NoError(t, err)

	validSession := data.Session{
		IsActive:      true,
		TransactionId: "1234",
		TagId:         "1234",
		Started:       time.Now().Format(time.RFC3339),
		Consumption:   nil,
	}

	expiredSession := data.Session{
		IsActive:      true,
		TransactionId: "1234",
		TagId:         "1234",
		Started:       time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local).Format(time.RFC3339),
		Consumption:   nil,
	}

	invalidSession := data.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	// ok case
	connector.SetStatus(core.ChargePointStatusCharging, core.NoError)
	err, timeElapsed := connector.ResumeCharging(validSession)
	require.NoError(err)
	require.InDelta(0, timeElapsed, 1)

	err = connector.StopCharging(core.ReasonLocal)
	require.NoError(err)

	//invalid session
	connector.SetStatus(core.ChargePointStatusCharging, core.NoError)
	err, timeElapsed = connector.ResumeCharging(invalidSession)
	require.Error(err)
	require.InDelta(connector.MaxChargingTime, timeElapsed, 1)

	//invalid status
	connector.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	err, timeElapsed = connector.ResumeCharging(validSession)
	require.Error(err)
	require.InDelta(connector.MaxChargingTime, timeElapsed, 1)

	//expired session
	connector.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	err, timeElapsed = connector.ResumeCharging(expiredSession)
	require.Error(err)
	require.InDelta(connector.MaxChargingTime, timeElapsed, 1)

	// invalid status
	connector.SetStatus(core.ChargePointStatusUnavailable, core.EVCommunicationError)
	err, timeElapsed = connector.ResumeCharging(invalidSession)
	require.Error(err)
	require.InDelta(connector.MaxChargingTime, timeElapsed, 1)

}

func TestConnector_StartCharging(t *testing.T) {
	require := require.New(t)

	connector1, _ := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(14, false), nil, false, 15)

	//ok case
	err := connector1.StartCharging("1234", "1234")
	require.NoError(err)

	// cannot start new session on a connector that is already charging
	err = connector1.StartCharging("1234a", "1234a")
	require.Error(err)

	err = connector1.StopCharging(core.ReasonLocal)
	assert.NoError(t, err)

	//invalid transaction and tag id
	err = connector1.StartCharging("@1234asd", "~ˇˇ3123")
	require.Error(err)
}

func TestConnector_StopCharging(t *testing.T) {
	require := require.New(t)

	connector1, _ := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(16, false), nil, false, 15)

	//start charging
	err := connector1.StartCharging("1234", "1234")
	require.NoError(err)

	//ok case
	err = connector1.StopCharging(core.ReasonLocal)
	require.NoError(err)

	// cannot stop charging if the connector is available
	err = connector1.StopCharging(core.ReasonLocal)
	require.Error(err)
}

func TestNewConnector(t *testing.T) {
	require := require.New(t)

	// ok case
	connector1, err := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(21, false), nil, false, 15)
	require.NoError(err)

	require.Equal(1, connector1.ConnectorId)
	require.Equal(1, connector1.EvseId)
	require.Equal(core.ChargePointStatusAvailable, connector1.ConnectorStatus)
	require.Equal("Schuko", connector1.ConnectorType)
	require.Equal(15, connector1.MaxChargingTime)
	require.False(connector1.PowerMeterEnabled)

	//invalid evseId
	_, err = chargepoint.NewConnector(0, 1, "Schuko",
		hardware.NewRelay(21, false), nil, false, 15)
	require.Error(err)

	// invalid connectorId
	_, err = chargepoint.NewConnector(1, 0, "Schuko",
		hardware.NewRelay(12, false), nil, false, 15)
	require.Error(err)

	//negative connector id
	_, err = chargepoint.NewConnector(1, -1, "Schuko",
		hardware.NewRelay(13, false), nil, false, 15)
	require.Error(err)

	// relay already in use
	_, err = chargepoint.NewConnector(1, 5, "Schuko",
		hardware.NewRelay(21, false), nil, false, 15)
	require.Error(err)

}

func TestConnector_SamplePowerMeter(t *testing.T) {
	//todo
}

func TestConnector_ReserveConnector(t *testing.T) {
	require := require.New(t)

	connector1, err := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(25, false), nil, false, 15)
	require.NoError(err)

	//ok case
	err = connector1.ReserveConnector(1)
	require.NoError(err)

	// connector already reserved
	err = connector1.ReserveConnector(2)
	require.Error(err)

	err = connector1.RemoveReservation()
	assert.NoError(t, err)

	// invalid connector status
	connector1.SetStatus(core.ChargePointStatusCharging, core.NoError)
	err = connector1.ReserveConnector(2)
	require.Error(err)
}

func TestConnector_RemoveReservation(t *testing.T) {
	require := require.New(t)

	connector1, err := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(23, false), nil, false, 15)
	require.NoError(err)

	//make a reservation
	err = connector1.ReserveConnector(1)
	require.NoError(err)

	//ok case
	err = connector1.RemoveReservation()
	require.NoError(err)

	// cannot remove reservation that is not there
	err = connector1.RemoveReservation()
	require.Error(err)
}
