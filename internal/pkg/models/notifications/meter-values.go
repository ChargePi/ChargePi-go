package notifications

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type MeterValueNotification struct {
	ConnectorId   *int
	EvseId        int
	TransactionId *int
	MeterValues   []types.MeterValue
}

func NewMeterValueNotification(evseId int, connectorId, transactionId *int, meterValues ...types.MeterValue) MeterValueNotification {
	return MeterValueNotification{
		ConnectorId:   connectorId,
		EvseId:        evseId,
		TransactionId: transactionId,
		MeterValues:   meterValues,
	}
}
