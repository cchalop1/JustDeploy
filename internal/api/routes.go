package api

import (
	"cchalop1.com/deploy/internal/api/handlers"
	"cchalop1.com/deploy/internal/api/middleware"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"time"
)

func CreateRoutes(app *Application, deployService *service.DeployService) {
	// Create JWT auth middleware
	apiKeyAuth := middleware.JWTAuth(deployService)

	// Create request logger middleware
	// requestLogger := middleware.RequestLogger()

	// Apply middleware to all API routes
	api := app.Echo.Group("api")

	// Apply request logger to all API routes
	// api.Use(requestLogger)

	// Public endpoints (no API key required)
	api.POST("/github/events", handlers.PostGithubEvent(deployService))
	api.GET("/info", handlers.GetServerInfoHandler(deployService))
	authLimit := echomiddleware.RateLimiter(echomiddleware.NewRateLimiterMemoryStoreWithConfig(echomiddleware.RateLimiterMemoryStoreConfig{Rate: 0.1, Burst: 5, ExpiresIn: 10 * time.Minute}))
	api.POST("/setup", handlers.PostSaveInitialSetupHandler(deployService), authLimit)
	api.POST("/login", handlers.PostLoginHandler(deployService), authLimit)

	// Protected endpoints (API key required)
	api.Use(apiKeyAuth)
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, name := range []string{"id", "serviceId", "deployId"} {
				if value := c.Param(name); value != "" && !application.ValidID(value) {
					return echo.NewHTTPError(400, "invalid identifier")
				}
			}
			return next(c)
		}
	})
	api.POST("/logout", func(c echo.Context) error {
		if !deployService.Operation.TryLock() {
			return echo.NewHTTPError(409, "another operation is in progress")
		}
		defer deployService.Operation.Unlock()
		secret, err := application.GenerateJwtSecret()
		if err != nil {
			return err
		}
		if err := deployService.DatabaseAdapter.RotateJWTSecret(secret); err != nil {
			return err
		}
		return c.JSON(200, map[string]string{"message": "All sessions revoked"})
	})
	api.POST("/github/state", handlers.GithubStateHandler(deployService))

	// API routes
	api.GET("/deploy", handlers.GetDeployListHandler(deployService))
	api.GET("/server/:id/logs", handlers.GetServerProxyLogs(deployService))
	api.GET("/server/:id/deploy", handlers.GetDeployByServerIdHandler(deployService))

	api.POST("/deploy/config/", legacyUnsupported)
	api.POST("/deploy/config/:deployId", legacyUnsupported)

	api.PUT("/deploy/edit", legacyUnsupported)

	api.DELETE("/deploy/remove/:id", handlers.RemoveApplicationHandler(deployService))
	api.POST("/deploy/start/:id", handlers.StartAppHandler(deployService))
	api.POST("/deploy/stop/:id", handlers.StopAppHandler(deployService))
	api.POST("/deploy/redeploy/:id", handlers.ReDeployAppHandler(deployService))

	api.GET("/service/:productId", handlers.GetServicesListHandler(deployService))
	api.GET("/services", handlers.GetServicesHandler(deployService))
	api.GET("/deploy/:deployId/service-docker-compose", legacyUnsupported)

	// Service logs endpoints
	api.GET("/service/:serviceId/build-logs", handlers.GetServiceBuildLogs(deployService))
	api.GET("/service/:serviceId/run-logs", handlers.GetServiceRunLogs(deployService))
	api.GET("/service/:serviceId/commit-info", handlers.GetServiceCommitInfo(deployService))

	api.PUT("/service", handlers.UpdateServiceHandler(deployService))

	api.POST("/server/domain", handlers.PostAddDomainToServerById(deployService))
	api.PUT("/server/settings", handlers.ServerSettingsHandler(deployService))
	api.PUT("/server/tls-settings", handlers.PutTlsServerSettings(deployService))

	api.DELETE("/service/:serviceId", handlers.DeleteServiceHandler(deployService))

	// Github
	api.GET("/github/is-connected", handlers.GetGithubIsConnectedHandler(deployService))

	api.POST("/github/connect/:code", handlers.PostConnectGithubAppHandler(deployService))

	api.GET("/github/repos", handlers.GetGithubRepos(deployService))

	// Create Service
	api.POST("/database/create", handlers.PostCreateDatabaseHandler(deployService))
	api.POST("/repo/create", handlers.PostCreateRepoHandler(deployService))

	api.POST("/github/save-access-token/:installationId", handlers.PostSaveAccessTokenHandler(deployService))

	// Deploy
	api.POST("/deploy", handlers.PostDeployHandler(deployService))
}

func legacyUnsupported(c echo.Context) error {
	return echo.NewHTTPError(410, "This legacy endpoint is no longer supported; use the service API")
}
