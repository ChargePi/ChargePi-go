package chargepoint

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/http"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/evse"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	userDatabase "github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/service"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
)

// SetupApi Runs a gRPC API server at a specified address if it is enabled. The API is protected by an authentication layer.
// Check the user manual for defaults.
func SetupApi(
	db *badger.DB,
	api settings.Api,
	handler chargePoint.ChargePoint,
	tagManager auth.TagManager,
	manager evse.Manager,
	ocppVariableManager ocppConfigManager.Manager,
) {
	if !api.Enabled {
		log.Info("API is disabled")
		return
	}

	// User database layer
	userDb := userDatabase.NewUserDb(db)

	// User service layer
	userService := service.NewUserService(userDb)

	// Expose the API endpoints
	server := grpc.NewServer(api, handler, tagManager, manager, ocppVariableManager, userService)
	server.Run()
}

// SetupUi Runs a management UI server if enabled
func SetupUi(uiSettings settings.Ui) {
	if !uiSettings.Enabled {
		log.Info("Management UI is disabled")
		return
	}

	ui := http.NewUi()
	ui.Serve(uiSettings.Address)
}

// Creates a healthcheck endpoint
func setupHealthcheck() {
	log.Infof("Starting application healthcheck at localhost:8081")
	httpServer := http.NewAppServer()
	httpServer.Serve(":8081")
}
