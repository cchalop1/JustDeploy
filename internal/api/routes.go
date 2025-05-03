package api

import (
	"cchalop1.com/deploy/internal/api/handlers"
	"cchalop1.com/deploy/internal/api/middleware"
	"cchalop1.com/deploy/internal/api/service"
)

func CreateRoutes(app *Application, deployService *service.DeployService) {
	// Create API key auth middleware
	apiKeyAuth := middleware.APIKeyAuth(deployService)

	// Apply middleware to all API routes
	api := app.Echo.Group("api")

	// Public endpoints (no API key required)
	api.POST("/github/events", handlers.PostGithubEvent(deployService))
	api.GET("/info", handlers.GetServerInfoHandler(deployService))
	api.POST("/setup", handlers.PostSaveInitialSetupHandler(deployService))

	// Protected endpoints (API key required)
	api.Use(apiKeyAuth)

	// API routes
	api.GET("/deploy", handlers.GetDeployListHandler(deployService))
	api.GET("/server/:id/logs", handlers.GetServerProxyLogs(deployService))
	api.GET("/server/:id/deploy", handlers.GetDeployByServerIdHandler(deployService))

	api.POST("/deploy/config/", handlers.GetDeployConfigHandler(deployService))
	api.POST("/deploy/config/:deployId", handlers.GetDeployConfigHandler(deployService))

	api.PUT("/deploy/edit", handlers.EditDeployementHandler(deployService))
	// api.POST("/server", handlers.ConnectNewServer(deployService))

	api.DELETE("/deploy/remove/:id", handlers.RemoveApplicationHandler(deployService))
	api.POST("/deploy/start/:id", handlers.StartAppHandler(deployService))
	api.POST("/deploy/stop/:id", handlers.StopAppHandler(deployService))
	api.POST("/deploy/redeploy/:id", handlers.ReDeployAppHandler(deployService))

	api.GET("/service/:productId", handlers.GetServicesListHandler(deployService))
	api.GET("/services", handlers.GetServicesHandler(deployService))
	api.GET("/deploy/:deployId/service-docker-compose", handlers.GetServicesFromDockerComposeHandler(deployService))

	api.PUT("/service", handlers.UpdateServiceHandler(deployService))

	api.POST("/server/domain", handlers.PostAddDomainToServerById(deployService))
	api.PUT("/server/tls-settings", handlers.PutTlsServerSettings(deployService))

	api.DELETE("/service/:serviceId", handlers.DeleteServiceHandler(deployService))

	api.GET("/server/:id/loading", handlers.SubscriptionCreateServerLoadingState(deployService))
	api.GET("/deploy/:id/loading", handlers.SubscriptionCreateDeployLoadingState(deployService))

	// Github
	api.GET("/github/is-connected", handlers.GetGithubIsConnectedHandler(deployService))

	api.POST("/github/connect/:code", handlers.PostConnectGithubAppHandler(deployService))

	api.GET("/github/repos", handlers.GetGithubRepos(deployService))

	api.POST("/github/save-access-token/:installationId", handlers.PostSaveAccessTokenHandler(deployService))

	// Create Service
	api.POST("/database/create", handlers.PostCreateDatabaseHandler(deployService))
	api.POST("/repo/create", handlers.PostCreateRepoHandler(deployService))

	// Deploy
	api.POST("/deploy", handlers.PostDeployHandler(deployService))
}
