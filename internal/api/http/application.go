package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	ginlogrus "github.com/toorop/gin-logrus"
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
	u.router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())

	// Configure healthcheck
	err := healthcheck.New(u.router, config.DefaultConfig(), checks)
	if err != nil {
		log.WithError(err).Panic("Failed to configure healthcheck")
	}

	err = u.router.Run(url)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Fatal("Failed to start HTTP server")
	}
}
