package handlers

import (
	"fmt"
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func DeleteServiceHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serviceId := c.Param("serviceId")

		err := application.DeleteService(deployService, serviceId)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "error")
		}

		return c.JSON(http.StatusOK, true)
	}
}
