package application

import (
	"net/http"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type HttpAdapter struct {
	server          *echo.Echo
	deployConfig    domain.DeployConfigDto
	dockerAdapter   adapter.DockerAdapter
	databaseAdapter adapter.DatabaseAdapter
}

func NewHttpAdapter(deployConfig domain.DeployConfigDto) *HttpAdapter {
	HttpAdapter := HttpAdapter{
		deployConfig:  deployConfig,
		dockerAdapter: *adapter.NewDockerAdapter(),
	}

	if deployConfig.DeployStatus != "serverconfig" {
		HttpAdapter.dockerAdapter = ConnectToServer(deployConfig.ServerConfig)
	}
	return &HttpAdapter
}

func (http *HttpAdapter) createRoutes() {
	// Render web page
	http.server.Static("/", "web/dist")

	// api routes
	http.server.GET("/api/deploy", http.getFormDetails)
	http.server.POST("/api/deploy", http.postCreateDeployementRoute)
	http.server.POST("/api/connect", http.connectServerRoute)
	http.server.DELETE("/api/remove/:name", http.removeApplicationRoute)
	http.server.GET("/api/logs/:name", http.getApplicationLogsRoute)
	http.server.POST("/api/redeploy/:name", http.reDeployApplicationRoute)
}

func (http *HttpAdapter) StartServer(openBrowser bool) {
	http.server = echo.New()
	// Allow cors
	http.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	http.createRoutes()

	if openBrowser {
		adapter.OpenBrowser("http://localhost:8080")
	}
	http.server.Start(":8080")
}

func (h *HttpAdapter) getFormDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, h.deployConfig)
}

func (h *HttpAdapter) connectServerRoute(c echo.Context) error {
	connectServerDto := domain.ConnectServerDto{}
	err := c.Bind(&connectServerDto)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	h.dockerAdapter = ConnectToServer(connectServerDto)
	h.deployConfig.ServerConfig = connectServerDto
	h.deployConfig.DeployStatus = "appconfig"
	h.databaseAdapter.SaveState(h.deployConfig)

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Server is connected"})
}

func (h *HttpAdapter) postCreateDeployementRoute(c echo.Context) error {
	postCreateDeploymentRequest := domain.AppConfigDto{}
	err := c.Bind(&postCreateDeploymentRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	h.deployConfig.AppConfig = postCreateDeploymentRequest
	h.deployConfig.PathToProject = adapter.NewFilesystemAdapter().CleanPath(postCreateDeploymentRequest.PathToSource)

	deployService := NewDeploymentService(&h.dockerAdapter)
	deployService.DeployApplication(h.deployConfig)

	h.deployConfig.DeployStatus = "deployapp"
	if h.deployConfig.AppConfig.EnableTls {
		h.deployConfig.Url = "https://" + h.dockerAdapter.ServerDomain
	} else {
		h.deployConfig.Url = "http://" + h.dockerAdapter.ServerDomain
	}

	h.databaseAdapter.SaveState(h.deployConfig)
	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is deploy"})
}

func (h *HttpAdapter) removeApplicationRoute(c echo.Context) error {
	applicationName := c.Param("name")

	RemoveApplication(applicationName, &h.dockerAdapter)
	h.deployConfig.DeployStatus = "appconfig"

	h.databaseAdapter.SaveState(h.deployConfig)
	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is removed"})
}

func (h *HttpAdapter) getApplicationLogsRoute(c echo.Context) error {
	containerName := c.Param("name")

	logs := GetApplicationLogs(containerName, &h.dockerAdapter)
	return c.JSON(http.StatusOK, logs)
}

func (h *HttpAdapter) reDeployApplicationRoute(c echo.Context) error {
	// containerName := c.Param("name")
	// TODO: get the app with the container name when we mange multiple application

	deployService := NewDeploymentService(&h.dockerAdapter)

	h.dockerAdapter.Stop(h.deployConfig.AppConfig.Name)
	deployService.DeployApplication(h.deployConfig)

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is redeploy"})
}
