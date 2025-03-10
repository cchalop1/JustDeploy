package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Get all services
// @Description Retrieves all services managed by the application
// @Tags services
// @Accept json
// @Produce json
// @Success 200 {array} handlers.ServiceResponse "List of services"
// @Router /api/v1/services [get]
func GetServicesHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		services := application.GetServices(deployService)
		return c.JSON(http.StatusOK, services)
	}
}
