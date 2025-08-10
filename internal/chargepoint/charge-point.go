package chargepoint

import (
	"context"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/internal/evse/manager"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	v16 "github.com/ChargePi/ChargePi-go/internal/chargepoint/v16"
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	"github.com/ChargePi/ChargePi-go/internal/sessions/service/session"
	"github.com/ChargePi/ChargePi-go/pkg/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/models/ocpp"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
)

var supportedOcppV16Profiles = []string{
	core.ProfileName,
	reservation.ProfileName,
	remotetrigger.ProfileName,
	localauth.ProfileName,
}

// CreateChargePoint Creates a OCPP-enabled charge point based on the protocol version
func CreateChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger *zap.Logger,
	manager manager.Manager,
	tagManager auth.Manager,
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
		defaultOcppConfig, err := ocpp_v16.DefaultConfigurationFromProfiles(supportedOcppV16Profiles...)
		if err != nil {
			logger.With(zap.Error(err)).Fatal("Cannot create OCPP configuration")
		}

		ocppVariableManager, err := ocpp_v16.NewV16ConfigurationManager(*defaultOcppConfig, supportedOcppV16Profiles...)
		if err != nil {
			logger.With(zap.Error(err)).Fatal("Cannot create OCPP configuration manager")
		}

		// Create the OCPP 1.6 Charge Point
		err = settingsManager.SetOcppV16Manager(ocppVariableManager)
		if err != nil {
			logger.With(zap.Error(err)).Fatal("Cannot add OCPP configuration manager")
		}

		return v16.NewChargePoint(
			logger,
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
		logger.With(zap.String("protocolVersion", string(protocolVersion))).Fatal("Protocol version not supported")
		return nil
	}
}
