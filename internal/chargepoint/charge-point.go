package chargepoint

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/service/session"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
)

// CreateChargePoint Creates a OCPP-enabled charge point based on the protocol version
func CreateChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger *log.Logger,
	manager evse.Manager,
	tagManager auth.TagManager,
	sessionManager session.Manager,
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
		// Create the OCPP 1.6 Charge Point
		return v16.NewChargePoint(
			manager,
			tagManager,
			sessionManager,
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
