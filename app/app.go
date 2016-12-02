package app

import (
	"github.com/labstack/echo"
)

func RunApp(configFile string){
	e := echo.New()
	e.Start()
}