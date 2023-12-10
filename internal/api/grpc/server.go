package grpc

import (
	"context"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/service"
	grpc2 "github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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
	userService service.Service,
) *Server {
	var opts []grpc.ServerOption

	if settings.TLS.IsEnabled {
		// Add TLS if enabled
		tlsCredentials, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.PrivateKeyPath)
		if err != nil {
			log.WithError(err).Panic("Failed to fetch credentials")
		}

		opts = []grpc.ServerOption{grpc.Creds(tlsCredentials)}
	}

	// Add authentication middleware
	opts = append(opts, grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		grpcauth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
			log.Debug("Authenticating request")
			token, err := grpcauth.AuthFromMD(ctx, "basic")
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "no basic header found: %v", err)
			}

			if userService.CheckPassword(token, token) {
				return nil, status.Errorf(codes.Unauthenticated, "invalid auth credentials: %v", err)
			}

			return ctx, nil
		}),
		grpcrecovery.UnaryServerInterceptor(),
	)))

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

	log.Infof("Exposing API endpoints at %s", s.address)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.WithError(err).Panicf("Unable to listen to provided address: %s", s.address)
	}

	err = s.server.Serve(listener)
	if err != nil {
		log.WithError(err).Panic("Cannot expose API")
	}
}
