package application

import (
	"embed"
	"net/http"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

//go:generate sh copyfiles.sh
//go:embed dist
var webAssets embed.FS

func (h *HttpAdapter) createRoutes() {
	// Render web page
	// http.server.Static("/", "web/dist")

	h.server.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "dist", // because files are located in `web` directory in `webAssets` fs
		Filesystem: http.FS(webAssets),
	}))

	// api routes
	h.server.GET("/api/deploy", h.getFormDetails)
	h.server.POST("/api/deploy", h.postCreateDeployementRoute)
	h.server.POST("/api/connect", h.connectServerRoute)
	h.server.DELETE("/api/remove/:name", h.removeApplicationRoute)
	h.server.GET("/api/logs/:name", h.getApplicationLogsRoute)
	h.server.POST("/api/redeploy/:name", h.reDeployApplicationRoute)
	h.server.POST("/api/start/:name", h.startApplicationRoute)
	h.server.POST("/api/stop/:name", h.stopApplicationRoute)
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

	h.deployConfig.AppStatus = "runing"
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

	h.dockerAdapter.Delete(h.deployConfig.AppConfig.Name, false)
	deployService.DeployApplication(h.deployConfig)

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is redeploy"})
}

func (h *HttpAdapter) stopApplicationRoute(c echo.Context) error {
	containerName := c.Param("name")

	h.dockerAdapter.Stop(containerName)
	h.deployConfig.AppStatus = "Stoped"

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is stoped"})
}

func (h *HttpAdapter) startApplicationRoute(c echo.Context) error {
	containerName := c.Param("name")

	h.dockerAdapter.Start(containerName)
	h.deployConfig.AppStatus = "Runing"

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is stoped"})
}
