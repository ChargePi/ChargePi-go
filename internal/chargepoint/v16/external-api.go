package v16

import (
	"github.com/xBlaz3kx/ChargePi-go/internal/api"
)

// StartCharging Start charging on the first available Connector. If there is no available Connector, reject the request.
func (cp *ChargePoint) StartCharging(tagId string, connectorId int) (*api.StartTransactionResponse, error) {
	return nil, nil
}

// StopCharging Stop charging a connector
func (cp *ChargePoint) StopCharging(tagId string, connectorId int) (*api.StopTransactionResponse, error) {
	return nil, nil
}

// GetConnectorStatus Notify the central system about the connector's status and updates the LED indicator.
func (cp *ChargePoint) GetConnectorStatus(evseId, connectorId int) (*api.GetConnectorStatusResponse, error) {
	return nil, nil
}
