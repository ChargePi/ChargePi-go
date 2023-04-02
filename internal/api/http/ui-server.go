package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

type UI struct {
	router *gin.Engine
}

func NewUi() *UI {
	return &UI{
		router: gin.Default(),
	}
}

func (u *UI) Serve(url string, checks ...checks.Check) {
	// Configure healthcheck
	err := healthcheck.New(u.router, config.DefaultConfig(), checks)
	if err != nil {
		return
	}

	u.router.LoadHTMLGlob("templates/**/*.html")

	u.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	u.router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	err = u.router.Run(url)
	if err != nil {
		return
	}
}
