package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func PostConnectGithubAppHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.Param("code")
		if code == "" {
			return c.JSON(http.StatusBadRequest, "Code is required")
		}

		res, err := application.ConnectGithubApp(deployService, code)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, res)
	}
}
