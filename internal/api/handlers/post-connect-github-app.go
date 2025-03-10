package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Connect GitHub application
// @Description Connects the application to GitHub using an authorization code
// @Tags github
// @Accept json
// @Produce json
// @Param code path string true "GitHub authorization code"
// @Success 200 {object} interface{} "Connection result"
// @Failure 400 {string} string "Code is required"
// @Failure 500 {string} string "Error message"
// @Router /api/v1/github/connect/{code} [post]
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
