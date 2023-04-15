package chargepoint

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/http"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	userDatabase "github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/service"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
)

func SetupApi(
	db *badger.DB,
	api settings.Api,
	handler chargePoint.ChargePoint,
	tagManager auth.TagManager,
	manager evse.Manager,
	ocppVariableManager ocppConfigManager.Manager,
) {
	if !api.Enabled {
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

func SetupUi(uiSettings settings.Ui) {
	if !uiSettings.Enabled {
		return
	}

	ui := http.NewUi()
	ui.Serve(uiSettings.Address)
}

func setupApplicationHealthcheck() {
	log.Infof("Starting application healthcheck at localhost:8081")
	http.NewAppServer().Serve(":8081")
}
