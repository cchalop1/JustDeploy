package application

import (
	"net/http"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type HttpAdapter struct {
	server        *echo.Echo
	deployConfig  domain.DeployConfigDto
	dockerAdapter adapter.DockerAdapter
}

func NewHttpAdapter(deployConfig domain.DeployConfigDto) *HttpAdapter {
	return &HttpAdapter{
		deployConfig: deployConfig,
	}
}

func (http *HttpAdapter) createRoutes() {
	// Render web page
	http.server.Static("/", "web/dist")

	// api routes
	http.server.GET("/api/deploy", http.getFormDetails)
	http.server.POST("/api/deploy", http.postCreateDeployementRoute)
	http.server.POST("/api/connect", http.connectServerRoute)
	http.server.DELETE("/api/remove", http.removeApplicationRoute)
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
	h.deployConfig.DeployFromStatus = "appconfig"

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

	h.deployConfig.DeployFromStatus = "deployapp"
	if h.deployConfig.AppConfig.EnableTls {
		h.deployConfig.Url = "https://" + h.dockerAdapter.ServerDomain
	} else {
		h.deployConfig.Url = "http://" + h.dockerAdapter.ServerDomain
	}

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is deploy"})
}

func (h *HttpAdapter) removeApplicationRoute(c echo.Context) error {
	deployService := NewDeploymentService(&h.dockerAdapter)
	deployService.RemoveApplication(h.deployConfig.AppConfig.Name)
	h.deployConfig.DeployFromStatus = "appconfig"

	return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is removed"})
}
