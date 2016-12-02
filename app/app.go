package app

import (
	"github.com/allen13/kube-source/app/config"
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/echo/middleware"
)

func RunApp()(err error){
	config.Load()

	e := echo.New()
	BuildApp(e)

	bindAddr := config.Get("address") + ":" + config.Get("port")

	if config.Get("server")  == "https" {
		err = e.StartTLS(bindAddr, config.Get("tls_cert"), config.Get("tls_key"))
	} else {
		err = e.Start(bindAddr)
	}

	return
}

func BuildApp(e *echo.Echo){
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ok := map[string]string{"status":"ok"}
	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, ok)
	})
}