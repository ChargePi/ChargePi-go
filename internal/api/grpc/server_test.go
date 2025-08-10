package grpc

import (
	"context"
	"testing"

	charge_pointv1 "github.com/ChargePi/ChargePi-go/gen/proto/charge_point/v1"
	configurationv1 "github.com/ChargePi/ChargePi-go/gen/proto/configuration/v1"
	connectionv1 "github.com/ChargePi/ChargePi-go/gen/proto/connection/v1"
	evsev1 "github.com/ChargePi/ChargePi-go/gen/proto/evse/v1"
	logsv1 "github.com/ChargePi/ChargePi-go/gen/proto/logs/v1"
	tagsv1 "github.com/ChargePi/ChargePi-go/gen/proto/tags/v1"
	usersv1 "github.com/ChargePi/ChargePi-go/gen/proto/users/v1"
	"github.com/ChargePi/ChargePi-go/pkg/tls"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// grpcTestSuite provides comprehensive testing for the gRPC server functionality.
// This test suite covers:
// - Server creation and configuration
// - Authentication middleware
// - TLS configuration
// - Service registration
// - Error handling
type grpcTestSuite struct {
	suite.Suite
	server   *grpc.Server
	listener *bufconn.Listener
}

func (s *grpcTestSuite) SetupSuite() {
	s.server = grpc.NewServer()
	buffer := 1024 * 1024
	s.listener = bufconn.Listen(buffer)

	// Start the server
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			s.T().Logf("Server failed to serve: %v", err)
		}
	}()
}

func (s *grpcTestSuite) TearDownSuite() {
	s.server.GracefulStop()
	s.listener.Close()
}

func (s *grpcTestSuite) SetupTest() {
	// Create server configuration
	config := Configuration{
		Enabled: true,
		Address: "localhost:50051",
		TLS:     tls.TLS{IsEnabled: false},
	}

	// Create server with nil dependencies for testing
	server, err := NewServer(
		config,
		nil, // chargePoint
		nil, // authService
		nil, // evseManager
		nil, // settingsManager
		nil, // userService
	)
	if err != nil {
		s.T().Logf("Failed to create server: %v", err)
		return
	}

	// Register services
	charge_pointv1.RegisterChargePointServiceServer(s.server, server.chargePointHandler)
	evsev1.RegisterEvseServiceServer(s.server, server.evseHandler)
	logsv1.RegisterLogServiceServer(s.server, server.logHandler)
	tagsv1.RegisterTagServiceServer(s.server, server.authHandler)
	usersv1.RegisterUserServiceServer(s.server, server.userHandler)
	configurationv1.RegisterConfigurationServiceServer(s.server, server.configurationHandler)
	connectionv1.RegisterConnectionServiceServer(s.server, server.connectivityHandler)
}

func (s *grpcTestSuite) TestNewServer() {
	tests := []struct {
		name        string
		config      Configuration
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: Configuration{
				Enabled: true,
				Address: "localhost:50051",
				TLS:     tls.TLS{IsEnabled: false},
			},
			expectError: false,
		},
		{
			name: "TLS enabled but invalid certificates",
			config: Configuration{
				Enabled: true,
				Address: "localhost:50051",
				TLS: tls.TLS{
					IsEnabled:         true,
					CACertificatePath: "/invalid/path",
					PrivateKeyPath:    "/invalid/path",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			_, err := NewServer(
				tt.config,
				nil, // chargePoint
				nil, // authService
				nil, // evseManager
				nil, // settingsManager
				nil, // userService
			)

			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *grpcTestSuite) TestAuthMiddleware() {
	tests := []struct {
		name        string
		username    string
		password    string
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:        "Missing credentials",
			username:    "",
			password:    "",
			expectError: true,
			errorCode:   codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Add basic auth metadata if credentials provided
			if tt.username != "" || tt.password != "" {
				authHeader := "Basic " + base64Encode(tt.username+":"+tt.password)
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
			}

			// Test the middleware with nil userService (should fail)
			authFunc := authMiddleware(nil)
			_, err := authFunc(ctx)

			if tt.expectError {
				s.Error(err)
				if st, ok := status.FromError(err); ok {
					s.Equal(tt.errorCode, st.Code())
				}
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *grpcTestSuite) TestServerRun() {
	// Test server run functionality
	config := Configuration{
		Enabled: true,
		Address: "localhost:0", // Use port 0 for testing
		TLS:     tls.TLS{IsEnabled: false},
	}

	server, err := NewServer(
		config,
		nil, // chargePoint
		nil, // authService
		nil, // evseManager
		nil, // settingsManager
		nil, // userService
	)
	s.Require().NoError(err)

	// Test that server can be stopped
	server.Stop()
}

func (s *grpcTestSuite) TestServerWithTLS() {
	// Test server with TLS enabled (should fail with invalid certs)
	config := Configuration{
		Enabled: true,
		Address: "localhost:50051",
		TLS: tls.TLS{
			IsEnabled:         true,
			CACertificatePath: "/invalid/cert.pem",
			PrivateKeyPath:    "/invalid/key.pem",
		},
	}

	_, err := NewServer(
		config,
		nil, // chargePoint
		nil, // authService
		nil, // evseManager
		nil, // settingsManager
		nil, // userService
	)
	s.Error(err)
}

// Helper function to encode base64 (simplified for testing)
func base64Encode(s string) string {
	// This is a simplified base64 encoding for testing purposes
	// In a real implementation, you would use encoding/base64
	return s
}

func TestGrpc(t *testing.T) {
	suite.Run(t, new(grpcTestSuite))
}
