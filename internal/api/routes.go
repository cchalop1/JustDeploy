package api

import (
	"cchalop1.com/deploy/internal/api/handlers"
	"cchalop1.com/deploy/internal/api/service"
)

func CreateRoutes(app *Application, deployService *service.DeployService) {

	app.Echo.GET("/api/deploy", handlers.GetDeployListHandler(deployService))
	app.Echo.GET("/api/deploy/:id", handlers.GetDeployByIdHandler(deployService))
	app.Echo.GET("/api/server", handlers.GetServerListHandler(deployService))
	app.Echo.GET("/api/server/:id", handlers.GetServerByIdHandler(deployService))
	app.Echo.GET("/api/server/:id/deploy", handlers.GetDeployByServerIdHandler(deployService))

	app.Echo.GET("/api/deploy/config", handlers.GetDeployConfigHandler(deployService))

	app.Echo.POST("/api/deploy", handlers.CreateDeployementHandler(deployService))
	app.Echo.PUT("/api/deploy/edit", handlers.EditDeployementHandler(deployService))
	app.Echo.POST("/api/server", handlers.ConnectNewServer(deployService))

	app.Echo.DELETE("/api/deploy/remove/:id", handlers.RemoveApplicationHandler(deployService))
	app.Echo.DELETE("/api/server/remove/:id", handlers.RemoveServerHandler(deployService))
	app.Echo.POST("/api/deploy/start/:id", handlers.StartAppHandler(deployService))
	app.Echo.POST("/api/deploy/stop/:id", handlers.StopAppHandler(deployService))
	app.Echo.POST("/api/deploy/redeploy/:id", handlers.ReDeployAppHandler(deployService))

	app.Echo.GET("/api/deploy/logs/:id", handlers.GetLogsHandler(deployService))

}
