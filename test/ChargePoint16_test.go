package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"os/exec"
	"testing"
)

func TestChargePointHandler_AddConnectors(t *testing.T) {

}

func TestChargePointHandler_findAvailableConnector(t *testing.T) {
	assert := assert.New(t)
	connector1 := chargepoint.Connector{
		EvseId:          1,
		ConnectorId:     1,
		ConnectorType:   "Schuko",
		MaxChargingTime: 0,
		ConnectorStatus: core.ChargePointStatusAvailable,
	}
	connector2 := chargepoint.Connector{
		EvseId:          1,
		ConnectorId:     1,
		ConnectorType:   "Schuko",
		MaxChargingTime: 0,
		ConnectorStatus: core.ChargePointStatusAvailable,
	}
	connector3 := chargepoint.Connector{
		EvseId:          1,
		ConnectorId:     1,
		ConnectorType:   "Schuko",
		MaxChargingTime: 0,
		ConnectorStatus: core.ChargePointStatusAvailable,
	}
	handler := chargepoint.ChargePointHandler{
		Connectors: []*chargepoint.Connector{
			&connector1, &connector2, &connector3,
		},
	}

	assert.Equal(&connector1, handler.FindAvailableConnector())

	connector1.ConnectorStatus = core.ChargePointStatusUnavailable
	assert.Equal(&connector2, handler.FindAvailableConnector())

	connector2.ConnectorStatus = core.ChargePointStatusUnavailable
	connector3.ConnectorStatus = core.ChargePointStatusUnavailable

	// no available connectors
	assert.Nil(handler.FindAvailableConnector())

}

func TestChargePointHandler_findConnectorWithId(t *testing.T) {
	assert := assert.New(t)
	connector1 := chargepoint.Connector{
		EvseId:          1,
		ConnectorId:     1,
		ConnectorType:   "Schuko",
		MaxChargingTime: 0,
		ConnectorStatus: core.ChargePointStatusAvailable,
	}

	connector2 := chargepoint.Connector{
		EvseId:          1,
		ConnectorId:     1,
		ConnectorType:   "Schuko",
		MaxChargingTime: 0,
		ConnectorStatus: core.ChargePointStatusAvailable,
	}

	connector3 := chargepoint.Connector{
		EvseId:          1,
		ConnectorId:     1,
		ConnectorType:   "Schuko",
		MaxChargingTime: 0,
		ConnectorStatus: core.ChargePointStatusAvailable,
	}

	handler := chargepoint.ChargePointHandler{
		Connectors: []*chargepoint.Connector{
			&connector1, &connector2, &connector3,
		},
	}

	assert.Equal(&connector1, handler.FindConnectorWithId(1))

	assert.Nil(handler.FindConnectorWithId(5))

}

func TestChargePointHandler_findConnectorWithTagId(t *testing.T) {
	assert := assert.New(t)
	connector1, _ := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(15, false), nil, false, 0)

	connector2, _ := chargepoint.NewConnector(1, 2, "Schuko",
		hardware.NewRelay(25, false), nil, false, 0)

	connector3, _ := chargepoint.NewConnector(1, 3, "Schuko",
		hardware.NewRelay(31, false), nil, false, 0)

	handler := chargepoint.ChargePointHandler{
		Connectors: []*chargepoint.Connector{
			connector1, connector2, connector3,
		},
	}

	assert.NoError(connector2.StartCharging("test123", "test123"))

	conn := handler.FindConnectorWithTagId("test123")
	assert.Equal(connector2, conn)

	assert.Nil(handler.FindConnectorWithTagId("4"))
}

func TestChargePointHandler_findConnectorWithTransactionId(t *testing.T) {
	assert := assert.New(t)
	connector1, _ := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(11, false), nil, false, 0)

	connector2, _ := chargepoint.NewConnector(1, 2, "Schuko",
		hardware.NewRelay(22, false), nil, false, 0)

	connector3, _ := chargepoint.NewConnector(1, 3, "Schuko",
		hardware.NewRelay(33, false), nil, false, 0)

	handler := chargepoint.ChargePointHandler{
		Connectors: []*chargepoint.Connector{
			connector1, connector2, connector3,
		},
	}

	assert.NoError(connector2.StartCharging("test123", "test123"))

	assert.Equal(connector2, handler.FindConnectorWithTransactionId("test123"))

	assert.Nil(handler.FindConnectorWithTransactionId("transaction1234"))
}

func Test_getTLSClient(t *testing.T) {
	assert := assert.New(t)
	InvalidCACertificatePath, CACertificatePath, ClientCertificatePath, ClientKeyPath := "certs/cs/ca123.crt", "certs/ca.crt", "certs/cp/charge-point.crt", "certs/cp/charge-point.key"

	exec.Command("sudo ./create-test-certs.sh")

	// invalid paths
	assert.Nil(chargepoint.GetTLSClient(InvalidCACertificatePath, ClientCertificatePath, ClientKeyPath))
	assert.Nil(chargepoint.GetTLSClient(CACertificatePath, "certs/invalidCertificatePath.crt", ClientKeyPath))
	assert.Nil(chargepoint.GetTLSClient(CACertificatePath, ClientCertificatePath, "certs/cp/charge-point-invalid.key"))

	//valid combination
	assert.NotNil(chargepoint.GetTLSClient(CACertificatePath, ClientCertificatePath, ClientKeyPath))

}
