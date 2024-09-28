package chargepoint

import (
	"context"
	"errors"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/reader"
	"github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/ChargePi/ChargePi-go/pkg/tls"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	ocppDisplay "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
	"github.com/tavsec/gin-healthcheck/checks"
)

var FirmwareVersion = "develop"

var (
	ErrConnectorNil           = errors.New("connector pointer is nil")
	ErrChargePointUnavailable = errors.New("charge point unavailable")
	ErrTagUnauthorized        = errors.New("tag unauthorized")
)

type (
	Settings struct {
		// General information about the charge point
		Info Info `json:"info" yaml:"info" mapstructure:"info"`

		// Backend connection settings for the charge point
		ConnectionSettings ConnectionSettings `json:"connectionSettings" yaml:"connectionSettings" mapstructure:"connectionSettings"`

		// Common Hardware settings, including display, tag reader and indicator(s)
		Hardware Hardware `json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	}

	Hardware struct {
		Display display.Settings `json:"display" yaml:"display" mapstructure:"display"`

		TagReader reader.Settings `json:"reader" yaml:"reader" mapstructure:"reader"`

		Indicator indicator.Settings `json:"indicator" yaml:"indicator" mapstructure:"indicator"`
	}

	ConnectionSettings struct {
		// OCPP ID of the charge point
		Id string `json:"id,omitempty" yaml:"id" mapstructure:"id" validate:"required"`

		// OCPP version of the charge point
		ProtocolVersion ocpp.ProtocolVersion `json:"protocolVersion,omitempty" yaml:"protocolVersion" mapstructure:"protocolVersion" validate:"required"`

		// Connection URL, without the charge point ID. Must be a valid URL, e.g. ws://localhost:8080/ocpp/2.0.1
		ServerUri string `json:"uri,omitempty" yaml:"uri" mapstructure:"uri" validate:"required,url"`

		// Basic authentication details for the connection, optional
		BasicAuth *BasicAuth `json:"basicAuth,omitempty" yaml:"basicAuth" mapstructure:"basicAuth"`

		// TLS settings for the connection
		TLS tls.TLS `json:"tls" yaml:"tls" mapstructure:"tls"`
	}

	BasicAuth struct {
		// Username for basic authentication
		Username string `json:"user,omitempty" yaml:"user" mapstructure:"user"`

		// Password for basic authentication
		Password string `json:"password,omitempty" yaml:"password" mapstructure:"password"`
	}

	Info struct {
		// Maximum time allowed if free mode is enabled
		MaxChargingTime *int `json:"MaxChargingTime,omitempty" yaml:"MaxChargingTime" mapstructure:"MaxChargingTime"`

		// Whether the charge point free charging mode is enabled
		FreeMode bool `json:"freeMode,omitempty" yaml:"freeMode" mapstructure:"freeMode"`

		// AC or DC
		Type string `json:"type,omitempty" yaml:"type" mapstructure:"type" validate:"oneof=AC DC"`

		// in kW
		MaxPower float32 `json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`

		// OCPP details
		OCPPDetails OCPPDetails `json:"ocpp" yaml:"ocpp" mapstructure:"ocpp"`
	}

	FreeChargingMode struct {
		Enabled  bool   `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Strategy string `json:"strategy,omitempty" yaml:"strategy" mapstructure:"strategy"`
	}
)

// todo model, vendor & version should be baked in at compile time?
// todo pointSerialNumber serial number should be baked in at compile time or fetched from EVCC?
type OCPPDetails struct {
	Vendor                  string `json:"vendor" yaml:"vendor" mapstructure:"vendor" validate:"required"`
	Model                   string `json:"model" yaml:"model" mapstructure:"model" validate:"required"`
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" yaml:"chargeBoxSerialNumber,omitempty" mapstructure:"chargeBoxSerialNumber,omitempty"`
	ChargePointSerialNumber string `json:"pointSerialNumber" yaml:"pointSerialNumber" mapstructure:"pointSerialNumber"`
}

type ChargePoint interface {
	checks.Check

	// Lifecycle APIs
	Connect(ctx context.Context, serverUrl string) error
	Cleanup(reason core.Reason) error
	Reset(resetType string) error
	ApplyOpts(opts ...Options) error

	// Core functionality
	StartCharging(evseId, connectorId int, tagId string) error
	StopCharging(evseId, connectorId int, reason core.Reason) error
	StartChargingFreeMode(evseId int) error
	SetAvailability(availabilityType core.AvailabilityType) error

	GetVersion() string
	GetStatus() string
	IsConnected() bool

	// Connector APIs
	SendEVSEsDetails(evseId int, maxPower float32, connectors ...ocpp.Connector)
	ListenForConnectorStatusChange(ctx context.Context)

	// Options
	SetLogger(logger log.FieldLogger)

	// Display APIs
	SetDisplay(display display.Display) error
	DisplayMessage(display ocppDisplay.MessageInfo) error

	// Settings APIs
	SetIndicator(indicator indicator.Indicator) error
	SetIndicatorSettings(settings indicator.StatusMapping) error
	GetIndicatorSettings() indicator.StatusMapping

	// Reader
	SetReader(reader reader.Reader) error
	ListenForTag(ctx context.Context) (*string, error)

	// Settings
	SetSettings(settings Info) error
	GetSettings() Info

	SetConnectionSettings(settings ConnectionSettings) error
	GetConnectionSettings() ConnectionSettings
}
