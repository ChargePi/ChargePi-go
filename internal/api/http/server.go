package http

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"

	"github.com/ChargePi/ChargePi-go/pkg/tls"
)

// todo go:embed ui/build/*
var fs embed.FS

type Configuration struct {
	// Enabled is a flag to enable or disable the UI
	UiEnabled bool `json:"uiEnabled,omitempty" yaml:"uiEnabled" mapstructure:"uiEnabled"`

	// Address is the address where the UI will be served. It should be in the format of host:port
	Address string `json:"address,omitempty" yaml:"address" mapstructure:"address"`

	// TLS is the configuration for the TLS
	TLS tls.TLS `json:"tls,omitempty" yaml:"tls" mapstructure:"tls"`
}

type Server struct {
	router        *gin.Engine
	server        *http.Server
	logger        *zap.Logger
	configuration Configuration
}

func NewServer(configuration Configuration) *Server {
	ginRouter := gin.Default()
	logger := zap.L()

	// Add logging and recovery middleware
	ginRouter.Use(loggingMiddleware(logger), gin.Recovery())

	return &Server{
		router: ginRouter,
		server: &http.Server{
			Addr:    configuration.Address,
			Handler: ginRouter.Handler(),
		},
		logger: logger,
	}
}

func (u *Server) Serve(checks ...checks.Check) {
	// Configure healthcheck
	err := healthcheck.New(u.router, config.DefaultConfig(), checks)
	if err != nil {
		u.logger.With(zap.Error(err)).Panic("Failed to configure healthcheck")
	}

	// Setup UI
	if u.configuration.UiEnabled {
		u.logger.Sugar().Infof("Starting UI at %s", u.configuration.Address)
		u.router.StaticFS("/", http.FS(fs))
	}

	go func() {
		err = u.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			u.logger.With(zap.Error(err)).Fatal("Failed to start HTTP server")
		}
	}()
}

func (u *Server) Stop() {
	u.logger.Info("Stopping HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.server.Shutdown(ctx)
	if err != nil {
		u.logger.Fatal("Failed to stop HTTP server")
	}
}
