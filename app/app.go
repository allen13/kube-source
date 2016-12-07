package app

import (
	"github.com/allen13/kube-source/app/config"
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/echo/middleware"
)

func RunServer()(err error){
	config.Load()

	server := buildServer()

	bindAddr := config.Get("address") + ":" + config.Get("port")

	if config.Get("server")  == "https" {
		err = server.StartTLS(bindAddr, config.Get("tls_cert"), config.Get("tls_key"))
	} else {
		err = server.Start(bindAddr)
	}

	return
}

func buildServer()(e *echo.Echo){
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"status":"ok"})
	})

	return
}