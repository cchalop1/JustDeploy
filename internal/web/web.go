package web

import (
	"embed"
	"net/http"

	"cchalop1.com/deploy/internal/api"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed dist
var webAssets embed.FS

func CreateMiddlewareWebFiles(app *api.Application) {
	app.Echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "dist",
		Filesystem: http.FS(webAssets),
	}))

}
