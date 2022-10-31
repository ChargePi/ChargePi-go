package chargePoint

import (
	"context"
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
)

var (
	ErrConnectorNil               = errors.New("connector pointer is nil")
	ErrConnectorNotCharging       = errors.New("connector not charging")
	ErrNoConnectorWithTag         = errors.New("no connector with tag id")
	ErrNoConnectorWithTransaction = errors.New("no connector with transaction id")
	ErrNoAvailableConnectors      = errors.New("no available evses")
	ErrConnectorUnavailable       = errors.New("connector unavailable")
	ErrChargePointUnavailable     = errors.New("charge point unavailable")
	ErrTagUnauthorized            = errors.New("tag unauthorized")
)

type (
	ChargePoint interface {
		// Basic info
		Connect(ctx context.Context, serverUrl string)
		CleanUp(reason core.Reason)
		ApplyOpts(opts ...Options)

		// Core functionality
		StartCharging(evseId, connectorId int, tagId string) error
		StopCharging(evseId, connectorId int, reason core.Reason) error

		// Connector API
		AddEVSEs(evses []*settings.EVSE)
		ListenForConnectorStatusChange(ctx context.Context, ch <-chan StatusNotification)

		// Options
		SetLogger(logger *log.Logger)
		SetDisplay(display display.Display)
		SetIndicator(indicator indicator.Indicator)
		SetSettings(settings *settings.Settings)

		// Reader
		SetReader(reader reader.Reader)
		ListenForTag(ctx context.Context, tagChannel <-chan string) (*string, error)
	}
)
