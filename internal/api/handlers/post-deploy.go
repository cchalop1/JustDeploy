package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Deploy application
// @Description Deploys the current application using the deploy service
// @Tags deployment
// @Accept json
// @Produce json
// @Success 200 {object} dto.ResponseApi "Successfully deployed"
// @Failure 500 {object} dto.ResponseApi "Internal server error with error message"
// @Router /api/v1/deploy [post]
func PostDeployHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := application.DeployApplication(deployService)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Deployed"})

	}
}
