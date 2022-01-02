package chargepoint

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
)

type (
	ChargePoint interface {
		Init(ctx context.Context, settings *settings.Settings)
		Connect(ctx context.Context, serverUrl string)
		HandleChargingRequest(tagId string)
		CleanUp(reason core.Reason)
		ListenForTag(ctx context.Context, tagChannel <-chan string)
		AddConnectors(connectors []*settings.Connector)
		ListenForConnectorStatusChange(ctx context.Context, ch <-chan rxgo.Item)
	}
)
