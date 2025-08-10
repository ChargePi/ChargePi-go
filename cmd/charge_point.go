package cmd

import (
	"context"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	v16 "github.com/ChargePi/ChargePi-go/internal/chargepoint/v16"
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	settings "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/internal/sessions"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var supportedOcppV16Profiles = []string{
	core.ProfileName,
	reservation.ProfileName,
	remotetrigger.ProfileName,
	localauth.ProfileName,
}

// NewChargePoint Creates a OCPP-enabled charge point based on the protocol version
func NewChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger log.FieldLogger,
	manager manager.Manager,
	tagManager auth.Service,
	settingsManager settings.Manager,
	sessionManager sessions.Service,
	diagnosticsManager diagnostics.Service,
	hardware chargepoint.Hardware,
) (chargepoint.ChargePoint, error) {

	// Create a status indicator (if enabled)
	statusIndicator := indicator.NewIndicator(len(manager.GetEVSEs()), hardware.Indicator)

	// Attach additional components based on the configuration
	opts := []chargepoint.Options{
		chargepoint.WithDisplayFromSettings(hardware.Display),
		chargepoint.WithReaderFromSettings(ctx, hardware.TagReader),
		chargepoint.WithLogger(logger),
		chargepoint.WithIndicator(statusIndicator),
	}

	switch protocolVersion {
	case ocpp.OCPP16:
		// Setup OCPP configuration from the database
		return v16.NewChargePoint(manager, settingsManager, tagManager, sessionManager, diagnosticsManager, opts...)
	case ocpp.OCPP201:
		return nil, errors.New("Version 2.0.1 is not supported yet.")
	default:
		return nil, errors.New("Invalid protocol version")
	}
}
