package chargepoint

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/evse"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/diagnostics"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/service/session"
	"github.com/xBlaz3kx/ChargePi-go/pkg/indicator"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	"github.com/xBlaz3kx/ocppManager-go/ocpp_v16"
)

// CreateChargePoint Creates a OCPP-enabled charge point based on the protocol version
func CreateChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger log.FieldLogger,
	manager evse.Manager,
	tagManager auth.TagManager,
	sessionManager session.Manager,
	settingsManager cfg.Manager,
	diagnosticsManager diagnostics.Manager,
	hardware settings.Hardware,
) chargePoint.ChargePoint {

	// Create a status indicator (if enabled)
	statusIndicator := indicator.NewIndicator(len(manager.GetEVSEs()), hardware.Indicator)

	// Attach additional components based on the configuration
	opts := []chargePoint.Options{
		chargePoint.WithDisplayFromSettings(hardware.Display),
		chargePoint.WithReaderFromSettings(ctx, hardware.TagReader),
		chargePoint.WithLogger(logger),
		chargePoint.WithIndicator(statusIndicator),
	}

	switch protocolVersion {
	case ocpp.OCPP16:
		// Setup OCPP configuration from the database
		ocppVariableManager, err := ocpp_v16.NewV16ConfigurationManager(
			ocpp_v16.NewEmptyConfiguration(),
			core.ProfileName,
			reservation.ProfileName,
			remotetrigger.ProfileName,
			localauth.ProfileName,
		)

		// Create the OCPP 1.6 Charge Point
		err = settingsManager.SetOcppV16Manager(ocppVariableManager)
		if err != nil {
			logger.WithError(err).Fatal("Cannot add OCPP configuration manager")
		}

		return v16.NewChargePoint(
			manager,
			tagManager,
			sessionManager,
			diagnosticsManager,
			opts...,
		)
	case ocpp.OCPP201:
		logger.Fatal("Version 2.0.1 is not supported yet.")
		return nil
	default:
		logger.WithField("protocolVersion", protocolVersion).Fatal("Protocol version not supported")
		return nil
	}
}
