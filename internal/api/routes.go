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

	app.Echo.GET("/api/deploy/config", handlers.GetDeployConfigHandler(deployService))

	app.Echo.POST("/api/deploy", handlers.CreateDeployementHandler(deployService))
	app.Echo.POST("/api/server", handlers.ConnectNewServer(deployService))

	app.Echo.DELETE("/api/remove/:id", handlers.RemoveApplicationHandler(deployService))
	// app.Echo.GET("/api/logs/:name", handlers.GetLogsHandler(deployService))
	// app.Echo.POST("/api/redeploy/:name", handlers.ReDeployAppHandler(deployService))
	// app.Echo.POST("/api/start/:name", handlers.StartAppHandler(deployService))
	// app.Echo.POST("/api/stop/:name", handlers.StopAppHandler(deployService))
	// app.Echo.PUT("/api/deploy/:name", handlers.EditDeployementHandler(deployService))
}
