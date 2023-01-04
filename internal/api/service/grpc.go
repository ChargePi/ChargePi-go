package service

import (
	"net"

	log "github.com/sirupsen/logrus"
	grpc2 "github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/users"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"

	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	server             *grpc.Server
	address            string
	service            *Service
	authService        *AuthService
	chargePointService *ChargePointService
	logService         *LogService
	userService        *UserService
}

func NewServer(
	settings settings.Api,
	point chargePoint.ChargePoint,
	authCache auth.TagManager,
	manager evse.Manager,
	configurationManager ocppConfigManager.Manager,
	userService users.Service,
) *Server {
	var opts []grpc.ServerOption

	if settings.TLS.IsEnabled {
		creds, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.PrivateKeyPath)
		if err != nil {
			log.WithError(err).Panic("Failed to fetch credentials")
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	return &Server{
		server:             grpc.NewServer(opts...),
		address:            settings.Address,
		service:            NewEvseService(manager),
		authService:        NewAuthService(authCache),
		chargePointService: NewChargePointService(point, configurationManager),
		logService:         NewLogService(),
		userService:        NewUserService(userService),
	}
}

func (s *Server) Run() {
	grpc2.RegisterChargePointServer(s.server, s.chargePointService)
	grpc2.RegisterEvseServer(s.server, s.service)
	grpc2.RegisterLogServer(s.server, s.logService)
	grpc2.RegisterTagServer(s.server, s.authService)
	grpc2.RegisterUsersServer(s.server, s.userService)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.WithError(err).Panicf("Unable to listen to provided address: %s", s.address)
	}

	err = s.server.Serve(listener)
	if err != nil {
		log.WithError(err).Panic("Cannot expose API")
	}
}
