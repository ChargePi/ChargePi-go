package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mandrigin/gin-spa/spa"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type UI struct {
	router *gin.Engine
}

func NewUi() *UI {
	return &UI{
		router: gin.Default(),
	}
}

func (u *UI) Serve(url string) {
	log.Infof("Starting UI at %s", url)
	u.router.Use(ginlogrus.Logger(log.StandardLogger()), spa.Middleware("/", "./ui/build"))

	err := u.router.Run(url)
	if err != nil {
		return
	}
}
