package http

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

type App struct {
	router *gin.Engine
}

func NewAppServer() *App {
	return &App{
		router: gin.Default(),
	}
}

func (u *App) Serve(url string, checks ...checks.Check) {
	// Configure healthcheck
	err := healthcheck.New(u.router, config.DefaultConfig(), checks)
	if err != nil {
		log.WithError(err).Panic("Failed to configure healthcheck")
	}

	err = u.router.Run(url)
	if err != nil {
		return
	}
}
