package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Echo *echo.Echo
}

func NewApplication() *Application {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	return &Application{
		Echo: e,
	}
}

func (app *Application) StartServer() {
	app.Echo.Start(":8080")
}
