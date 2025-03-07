package api

import (
	"cchalop1.com/deploy/internal/api/handlers"
	"cchalop1.com/deploy/internal/api/service"
)

func CreateRoutes(app *Application, deployService *service.DeployService) {

	app.Echo.GET("/api/deploy", handlers.GetDeployListHandler(deployService))
	app.Echo.GET("/api/server/:id/logs", handlers.GetServerProxyLogs(deployService))
	app.Echo.GET("/api/server/:id/deploy", handlers.GetDeployByServerIdHandler(deployService))
	app.Echo.GET("/api/server/info", handlers.GetServerInfoHandler(deployService))

	app.Echo.POST("/api/deploy/config/", handlers.GetDeployConfigHandler(deployService))
	app.Echo.POST("/api/deploy/config/:deployId", handlers.GetDeployConfigHandler(deployService))

	app.Echo.PUT("/api/deploy/edit", handlers.EditDeployementHandler(deployService))
	app.Echo.POST("/api/server", handlers.ConnectNewServer(deployService))

	app.Echo.DELETE("/api/deploy/remove/:id", handlers.RemoveApplicationHandler(deployService))
	app.Echo.POST("/api/deploy/start/:id", handlers.StartAppHandler(deployService))
	app.Echo.POST("/api/deploy/stop/:id", handlers.StopAppHandler(deployService))
	app.Echo.POST("/api/deploy/redeploy/:id", handlers.ReDeployAppHandler(deployService))

	app.Echo.GET("/api/service/:productId", handlers.GetServicesListHandler(deployService))
	app.Echo.GET("/api/services", handlers.GetServicesHandler(deployService))
	app.Echo.GET("/api/deploy/:deployId/service-docker-compose", handlers.GetServicesFromDockerComposeHandler(deployService))

	app.Echo.PUT("/api/project/:projectId/service", handlers.UpdateServiceHandler(deployService))

	app.Echo.POST("/api/server/domain", handlers.PostAddDomainToServerById(deployService))

	app.Echo.DELETE("/api/service/:serviceId", handlers.DeleteServiceHandler(deployService))

	app.Echo.GET("/api/server/:id/loading", handlers.SubscriptionCreateServerLoadingState(deployService))
	app.Echo.GET("/api/deploy/:id/loading", handlers.SubscriptionCreateDeployLoadingState(deployService))

	app.Echo.GET("/api/version", handlers.GetVersionHandler(deployService))

	// Github
	app.Echo.GET("/api/github/is-connected", handlers.GetGithubIsConnectedHandler(deployService))

	app.Echo.POST("/api/github/connect/:code", handlers.PostConnectGithubAppHandler(deployService))

	app.Echo.GET("/api/github/repos", handlers.GetGithubRepos(deployService))

	app.Echo.POST("/api/github/save-access-token/:installationId", handlers.PostSaveAccessTokenHandler(deployService))

	// Create Service
	// app.Echo.POST("/api/database/create", handlers.PostCreateDatabaseHandler(deployService))
	app.Echo.POST("/api/repo/create", handlers.PostCreateRepoHandler(deployService))

	// Deploy
	app.Echo.POST("/api/deploy", handlers.PostDeployHandler(deployService))
}
