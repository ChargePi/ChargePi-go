package chargePoint

import (
	"context"
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	ocppDisplay "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/display"
	"github.com/xBlaz3kx/ChargePi-go/pkg/indicator"
	data "github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	settings2 "github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/reader"
)

var (
	ErrConnectorNil           = errors.New("connector pointer is nil")
	ErrChargePointUnavailable = errors.New("charge point unavailable")
	ErrTagUnauthorized        = errors.New("tag unauthorized")
)

type ChargePoint interface {
	// Lifecycle APIs
	Connect(ctx context.Context, serverUrl string)
	CleanUp(reason core.Reason)
	Reset(resetType string) error
	ApplyOpts(opts ...Options)

	// Core functionality
	StartCharging(evseId, connectorId int, tagId string) error
	StopCharging(evseId, connectorId int, reason core.Reason) error
	StartChargingFreeMode(evseId int) error

	// Connector APIs
	SendEVSEsDetails(evseId int, maxPower float32, connectors ...data.Connector)
	ListenForConnectorStatusChange(ctx context.Context, ch <-chan notifications.StatusNotification)

	// Options
	SetLogger(logger log.FieldLogger)

	// Display APIs
	SetDisplay(display display.Display) error
	DisplayMessage(display ocppDisplay.MessageInfo) error

	// Indicator APIs
	SetIndicator(indicator indicator.Indicator) error
	SetIndicatorSettings(settings settings2.IndicatorStatusMapping) error
	GetIndicatorSettings() settings2.IndicatorStatusMapping

	// Reader
	SetReader(reader reader.Reader) error
	ListenForTag(ctx context.Context, tagChannel <-chan string) (*string, error)

	// Settings
	SetSettings(settings settings.Info) error
	GetSettings() settings.Info

	// Connection settings
	SetConnectionSettings(settings settings.ConnectionSettings) error
	GetConnectionSettings() settings.ConnectionSettings

	GetVersion() string
	GetStatus() string
	// SetStatus(status string) error
	IsConnected() bool
}
