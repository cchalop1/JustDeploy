package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// GetServiceRunLogs returns a handler function for retrieving run logs for a service
func GetServiceRunLogs(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serviceId := c.Param("serviceId")
		if serviceId == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Service ID is required"})
		}

		logs, err := application.GetServiceRunLogs(deployService, serviceId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, logs)
	}
}
