package grpc

import (
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/users"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"

	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server struct {
	server  *grpc.Server
	address string
	service *Service
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
		creds, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.ClientKeyPath)
		if err != nil {
			log.WithError(err).Panic("Failed to fetch credentials")
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	return &Server{
		server:  grpc.NewServer(opts...),
		address: settings.Address,
		service: NewGrpcService(point, authCache, manager, configurationManager, userService),
	}
}

func (s *Server) Run() {
	RegisterChargePointServer(s.server, s.service)
	RegisterEvseServer(s.server, s.service)
	RegisterLogServer(s.server, s.service)
	RegisterTagServer(s.server, s.service)
	RegisterUserServer(s.server, s.service)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.WithError(err).Panicf("Unable to listen to provided address: %s", s.address)
	}

	err = s.server.Serve(listener)
	if err != nil {
		log.WithError(err).Panic("Cannot expose API")
	}
}
