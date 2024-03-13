package api

import (
	"cchalop1.com/deploy/internal/api/handlers"
	"cchalop1.com/deploy/internal/api/usecase"
)

func CreateRoutes(app *Application, deployUseCase *usecase.DeployUseCase) {
	app.Echo.GET("/api/deploy", handlers.GetDeployConfigHandler(deployUseCase))
	app.Echo.POST("/api/deploy", handlers.CreateDeployementHandler(deployUseCase))
	app.Echo.POST("/api/connect", handlers.ConnectAndSetupServerHandler(deployUseCase))
	app.Echo.DELETE("/api/remove/:name", handlers.RemoveApplicationHandler(deployUseCase))
	app.Echo.GET("/api/logs/:name", handlers.GetLogsHandler(deployUseCase))
	app.Echo.POST("/api/redeploy/:name", handlers.ReDeployAppHandler(deployUseCase))
	app.Echo.POST("/api/start/:name", handlers.StartAppHandler(deployUseCase))
	app.Echo.POST("/api/stop/:name", handlers.StopAppHandler(deployUseCase))
}
