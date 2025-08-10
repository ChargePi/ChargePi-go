package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"go.uber.org/zap"
)

type App struct {
	router *gin.Engine
}

func NewAppServer() *App {
	ginRouter := gin.Default()

	return &App{
		router: ginRouter,
	}
}

func (u *App) Serve(url string, checks ...checks.Check) {
	u.router.Use(loggingMiddleware(), gin.Recovery())

	// Configure healthcheck
	err := healthcheck.New(u.router, config.DefaultConfig(), checks)
	if err != nil {
		zap.L().Panic("Failed to configure healthcheck", zap.Error(err))
	}

	err = u.router.Run(url)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		zap.L().Fatal("Failed to start HTTP server", zap.Error(err))
	}
}
