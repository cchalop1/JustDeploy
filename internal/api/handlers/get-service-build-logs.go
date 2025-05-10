package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// GetServiceBuildLogs returns a handler function for retrieving build logs for a service
func GetServiceBuildLogs(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serviceId := c.Param("serviceId")
		if serviceId == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Service ID is required"})
		}

		logs, err := application.GetServiceBuildLogs(deployService, serviceId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, logs)
	}
}
