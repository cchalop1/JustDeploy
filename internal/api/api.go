package api

import (
	"cchalop1.com/deploy/internal/adapter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Echo *echo.Echo
}

func NewApplication() *Application {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	return &Application{
		Echo: e,
	}
}

func (app *Application) StartServer(openBrowser bool) {
	if openBrowser {
		adapter.OpenBrowser("http://localhost:8080")
	}
	app.Echo.Start(":8080")
}
