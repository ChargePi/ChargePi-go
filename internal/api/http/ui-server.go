package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UI struct {
}

func NewUi() *UI {
	return &UI{}
}

func (u *UI) Serve(url string) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	err := r.Run(url)
	if err != nil {
		return
	}
}
