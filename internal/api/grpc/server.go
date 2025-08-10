package grpc

import (
	"context"
	"fmt"
	"net"

	charge_pointv1 "github.com/ChargePi/ChargePi-go/gen/proto/charge_point/v1"
	configurationv1 "github.com/ChargePi/ChargePi-go/gen/proto/configuration/v1"
	connectionv1 "github.com/ChargePi/ChargePi-go/gen/proto/connection/v1"
	evsev1 "github.com/ChargePi/ChargePi-go/gen/proto/evse/v1"
	logsv1 "github.com/ChargePi/ChargePi-go/gen/proto/logs/v1"
	tagsv1 "github.com/ChargePi/ChargePi-go/gen/proto/tags/v1"
	usersv1 "github.com/ChargePi/ChargePi-go/gen/proto/users/v1"
	"github.com/ChargePi/ChargePi-go/internal/auth"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	settings "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/internal/users"
	"github.com/ChargePi/ChargePi-go/pkg/tls"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// GRPC API configuration
type Configuration struct {
	// Enabled is a flag to enable or disable the API
	Enabled bool `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`

	// Address is the address where the API will be served. It should be in the format of host:port.
	Address string `json:"address,omitempty" yaml:"address" mapstructure:"address"`

	// TLS is the configuration for the TLS
	TLS tls.TLS `json:"tls,omitempty" yaml:"tls" mapstructure:"tls"`
}

type Server struct {
	logger               *zap.Logger
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
	settings Configuration,
	point chargePoint.ChargePoint,
	authCache auth.Service,
	manager manager.Manager,
	settingsManager settings.Manager,
	userService users.Service,
) (*Server, error) {
	var opts []grpc.ServerOption

	logger := zap.L().Named("grpc_server")

	if settings.TLS.IsEnabled {
		// Add TLS if enabled
		tlsCredentials, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to TLS certificates")
		}

		opts = append(opts, grpc.Creds(tlsCredentials))
	}

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	// Add authentication, recovery and logging middleware
	opts = append(opts, grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		logging.UnaryServerInterceptor(interceptorLogger(logger), logOpts...),
		grpcauth.UnaryServerInterceptor(authMiddleware(logger,userService)),
		grpcrecovery.UnaryServerInterceptor(),
	)))

	return &Server{
		logger:               logger,
		server:               grpc.NewServer(opts...),
		address:              settings.Address,
		evseHandler:          NewEvseHandler(manager),
		authHandler:          NewAuthService(authCache),
		chargePointHandler:   NewChargePointService(point),
		logHandler:           NewLogHandler(),
		userHandler:          NewUserHandler(userService),
		configurationHandler: NewConfigurationHandler(settingsManager),
		connectivityHandler:  NewConnectivityHandler(),
	}, nil
}

func (s *Server) Run() error {
	charge_pointv1.RegisterChargePointServiceServer(s.server, s.chargePointHandler)
	evsev1.RegisterEvseServiceServer(s.server, s.evseHandler)
	logsv1.RegisterLogServiceServer(s.server, s.logHandler)
	tagsv1.RegisterTagServiceServer(s.server, s.authHandler)
	usersv1.RegisterUserServiceServer(s.server, s.userHandler)
	configurationv1.RegisterConfigurationServiceServer(s.server, s.configurationHandler)
	connectionv1.RegisterConnectionServiceServer(s.server, s.connectivityHandler)

	s.logger.Info("Exposing API endpoints", zap.String("address", s.address))

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("unable to listen to provided address: %s", s.address)
	}

	err = s.server.Serve(listener)
	if err != nil {
		return fmt.Errorf("unable to listen to provided address: %s", s.address)
	}

	return nil
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	s.server.GracefulStop()
}

// authMiddleware is a middleware function that authenticates incoming requests using basic auth.
func authMiddleware(logger *zap.Logger, userService users.Service) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		zap.L().Debug("Authenticating request")
		token, err := grpcauth.AuthFromMD(ctx, "basic")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "no basic header found: %v", err)
		}

		if userService.CheckPassword(token, token) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth credentials: %v", err)
		}

		return ctx, nil
	}
}

func interceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		for i.Next() {
			k, v := i.At()
			f[k] = v
		}
		l := l.WithFields(f)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
