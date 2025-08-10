package grpc

import (
	"context"
	"net"

	charge_pointv1 "github.com/ChargePi/ChargePi-go/gen/proto/charge_point/v1"
	configurationv1 "github.com/ChargePi/ChargePi-go/gen/proto/configuration/v1"
	connectionv1 "github.com/ChargePi/ChargePi-go/gen/proto/connection/v1"
	evsev1 "github.com/ChargePi/ChargePi-go/gen/proto/evse/v1"
	logsv1 "github.com/ChargePi/ChargePi-go/gen/proto/logs/v1"
	tagsv1 "github.com/ChargePi/ChargePi-go/gen/proto/tags/v1"
	usersv1 "github.com/ChargePi/ChargePi-go/gen/proto/users/v1"
	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	"github.com/ChargePi/ChargePi-go/internal/users/service"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type Server struct {
	server               *grpc.Server
	address              string
	evseHandler          *EvseHandler
	authHandler          *AuthService
	chargePointHandler   *ChargePointHandler
	logHandler           *LogHandler
	userHandler          *UserHandler
	configurationHandler *ConfigurationHandler
	connectivityHandler  *ConnectivityHandler
}

func NewServer(
	settings settings.Api,
	point chargePoint.ChargePoint,
	authCache auth.Manager,
	manager manager.Manager,
	settingsManager cfg.Manager,
	userService service.Service,
) *Server {
	var opts []grpc.ServerOption

	if settings.TLS.IsEnabled {
		// Add TLS if enabled
		tlsCredentials, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.PrivateKeyPath)
		if err != nil {
			zap.L().Panic("Failed to fetch credentials", zap.Error(err))
		}

		opts = []grpc.ServerOption{grpc.Creds(tlsCredentials)}
	}

	// Add authentication middleware
	opts = append(opts, grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		grpcauth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
			zap.L().Debug("Authenticating request")
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
		server:               grpc.NewServer(opts...),
		address:              settings.Address,
		evseHandler:          NewEvseHandler(manager),
		authHandler:          NewAuthService(authCache),
		chargePointHandler:   NewChargePointService(point, settingsManager),
		logHandler:           NewLogHandler(),
		userHandler:          NewUserHandler(userService),
		configurationHandler: NewConfigurationHandler(settingsManager),
		connectivityHandler:  NewConnectivityHandler(),
	}
}

func (s *Server) Run() {
	charge_pointv1.RegisterChargePointServiceServer(s.server, s.chargePointHandler)
	evsev1.RegisterEvseServiceServer(s.server, s.evseHandler)
	logsv1.RegisterLogServiceServer(s.server, s.logHandler)
	tagsv1.RegisterTagServiceServer(s.server, s.authHandler)
	usersv1.RegisterUserServiceServer(s.server, s.userHandler)
	configurationv1.RegisterConfigurationServiceServer(s.server, s.configurationHandler)
	connectionv1.RegisterConnectionServiceServer(s.server, s.connectivityHandler)

	zap.L().Info("Exposing API endpoints", zap.String("address", s.address))

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		zap.L().Panic("Unable to listen to provided address", zap.Error(err), zap.String("address", s.address))
	}

	err = s.server.Serve(listener)
	if err != nil {
		zap.L().Panic("Cannot expose API", zap.Error(err))
	}
}
