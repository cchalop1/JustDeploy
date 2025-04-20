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
	api := app.Echo.Group("")
	api.Use(apiKeyAuth)

	// API routes
	api.GET("/api/deploy", handlers.GetDeployListHandler(deployService))
	api.GET("/api/server/:id/logs", handlers.GetServerProxyLogs(deployService))
	api.GET("/api/server/:id/deploy", handlers.GetDeployByServerIdHandler(deployService))
	api.GET("/api/server/info", handlers.GetServerInfoHandler(deployService))

	api.POST("/api/deploy/config/", handlers.GetDeployConfigHandler(deployService))
	api.POST("/api/deploy/config/:deployId", handlers.GetDeployConfigHandler(deployService))

	api.PUT("/api/deploy/edit", handlers.EditDeployementHandler(deployService))
	api.POST("/api/server", handlers.ConnectNewServer(deployService))

	api.DELETE("/api/deploy/remove/:id", handlers.RemoveApplicationHandler(deployService))
	api.POST("/api/deploy/start/:id", handlers.StartAppHandler(deployService))
	api.POST("/api/deploy/stop/:id", handlers.StopAppHandler(deployService))
	api.POST("/api/deploy/redeploy/:id", handlers.ReDeployAppHandler(deployService))

	api.GET("/api/service/:productId", handlers.GetServicesListHandler(deployService))
	api.GET("/api/services", handlers.GetServicesHandler(deployService))
	api.GET("/api/deploy/:deployId/service-docker-compose", handlers.GetServicesFromDockerComposeHandler(deployService))

	api.PUT("/api/service", handlers.UpdateServiceHandler(deployService))

	api.POST("/api/server/domain", handlers.PostAddDomainToServerById(deployService))

	api.DELETE("/api/service/:serviceId", handlers.DeleteServiceHandler(deployService))

	api.GET("/api/server/:id/loading", handlers.SubscriptionCreateServerLoadingState(deployService))
	api.GET("/api/deploy/:id/loading", handlers.SubscriptionCreateDeployLoadingState(deployService))

	api.GET("/api/version", handlers.GetVersionHandler(deployService))

	// Github
	api.GET("/api/github/is-connected", handlers.GetGithubIsConnectedHandler(deployService))

	api.POST("/api/github/connect/:code", handlers.PostConnectGithubAppHandler(deployService))

	api.GET("/api/github/repos", handlers.GetGithubRepos(deployService))

	api.POST("/api/github/save-access-token/:installationId", handlers.PostSaveAccessTokenHandler(deployService))

	// Create Service
	// api.POST("/api/database/create", handlers.PostCreateDatabaseHandler(deployService))
	api.POST("/api/repo/create", handlers.PostCreateRepoHandler(deployService))

	// Deploy
	api.POST("/api/deploy", handlers.PostDeployHandler(deployService))
}
