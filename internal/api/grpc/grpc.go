package grpc

import (
	log "github.com/sirupsen/logrus"

	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server struct {
	server  *grpc.Server
	address string
	service *Service
}

func NewServer(settings settings.Api, point chargePoint.ChargePoint, authCache auth.TagManager) *Server {
	var opts []grpc.ServerOption

	if settings.TLS.IsEnabled {
		creds, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.ClientKeyPath)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	return &Server{
		server:  grpc.NewServer(opts...),
		address: settings.Address,
		service: NewChargePointService(point, authCache),
	}
}

func (s *Server) Run() {
	RegisterChargePointServer(s.server, s.service)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.WithError(err).Fatalf("Unable to listen to provided address: %s", s.address)
	}

	err = s.server.Serve(listener)
	if err != nil {
		log.WithError(err).Fatal("Cannot expose API")
	}
}
