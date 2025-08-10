package chargepoint

import (
	"github.com/ChargePi/ChargePi-go/internal/api/grpc"
	"github.com/ChargePi/ChargePi-go/internal/api/http"
	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	userDatabase "github.com/ChargePi/ChargePi-go/internal/users/pkg/database"
	"github.com/ChargePi/ChargePi-go/internal/users/service"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// SetupApi Runs a gRPC API server at a specified address if it is enabled. The API is protected by an authentication layer.
// Check the user manual for defaults.
func SetupApi(
	db *badger.DB,
	api settings.Api,
	handler chargePoint.ChargePoint,
	tagManager auth.Manager,
	manager manager.Manager,
	settingsManager cfg.Manager,
) {
	if !api.Enabled {
		zap.L().Info("API is disabled")
		return
	}

	// User database layer
	userDb := userDatabase.NewUserDb(db)

	// User service layer
	userService := service.NewUserService(zap.L(), userDb)

	// Expose the API endpoints
	server := grpc.NewServer(api, handler, tagManager, manager, settingsManager, userService)
	server.Run()
}

// SetupUi Runs a management UI server if enabled
func SetupUi(uiSettings settings.Ui) {
	if !uiSettings.Enabled {
		zap.L().Info("Management UI is disabled")
		return
	}

	ui := http.NewUi()
	ui.Serve(uiSettings.Address)
}

// Creates a healthcheck endpoint
func setupHealthcheck() {
	zap.L().Info("Starting application healthcheck at localhost:8081")
	httpServer := http.NewAppServer()
	httpServer.Serve(":8081")
}
