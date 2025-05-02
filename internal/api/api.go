package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Echo *echo.Echo
}

func NewApplication(port string) *Application {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:" + port, "http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "x-api-key"},
	}))

	return &Application{
		Echo: e,
	}
}

func (app *Application) StartServer(port string) {
	app.Echo.Start(":" + port)
}
