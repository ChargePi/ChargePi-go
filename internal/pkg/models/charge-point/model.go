package chargePoint

import (
	"context"
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	data "github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"

	ocppDisplay "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
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
	SetLogger(logger *log.Logger)

	// Display APIs
	SetDisplay(display display.Display) error
	DisplayMessage(display ocppDisplay.MessageInfo) error

	// Indicator APIs
	SetIndicator(indicator indicator.Indicator) error
	SetIndicatorSettings(settings settings.IndicatorStatusMapping) error
	GetIndicatorSettings() settings.IndicatorStatusMapping

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
