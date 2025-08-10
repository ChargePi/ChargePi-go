package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mandrigin/gin-spa/spa"
	"go.uber.org/zap"
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
	zap.L().Info("Starting UI", zap.String("url", url))
	u.router.Use(loggingMiddleware(), spa.Middleware("/", "./ui/build"))

	err := u.router.Run(url)
	if err != nil {
		return
	}
}
