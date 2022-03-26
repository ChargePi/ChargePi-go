package chargePoint

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/internal/api"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
)

type (
	ChargePoint interface {
		Init(settings *settings.Settings)
		Connect(ctx context.Context, serverUrl string)
		HandleChargingRequest(tagId string) (*api.HandleChargingResponse, error)
		StartCharging(tagId string, connectorId int) (*api.StartTransactionResponse, error)
		StopCharging(tagId string, connectorId int) (*api.StopTransactionResponse, error)
		GetConnectorStatus(evseId, connectorId int) (*api.GetConnectorStatusResponse, error)
		CleanUp(reason core.Reason)
		ListenForTag(ctx context.Context, tagChannel <-chan string)
		AddConnectors(connectors []*settings.Connector)
		ListenForConnectorStatusChange(ctx context.Context, ch <-chan rxgo.Item)
	}
)
