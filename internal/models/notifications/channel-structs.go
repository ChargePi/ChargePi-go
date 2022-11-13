package notifications

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type (
	MeterValueNotification struct {
		ConnectorId   *int
		EvseId        int
		TransactionId *int
		MeterValues   []types.MeterValue
	}

	StatusNotification struct {
		EvseId    int
		Status    string
		ErrorCode string
	}
)

func NewMeterValueNotification(evseId int, connectorId, transactionId *int, meterValues ...types.MeterValue) MeterValueNotification {
	return MeterValueNotification{
		ConnectorId:   connectorId,
		EvseId:        evseId,
		TransactionId: transactionId,
		MeterValues:   meterValues,
	}
}

func NewStatusNotification(evseId int, status, errorCode string) StatusNotification {
	return StatusNotification{
		Status:    status,
		EvseId:    evseId,
		ErrorCode: errorCode,
	}
}
