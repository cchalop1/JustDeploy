package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/usecase"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/models"
	"github.com/labstack/echo/v4"
)

func CreateDeployementHandler(deployUseCase *usecase.DeployUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		postCreateDeploymentRequest := models.AppConfigDto{}
		err := c.Bind(&postCreateDeploymentRequest)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		deployUseCase.DeployConfig.AppConfig = postCreateDeploymentRequest
		deployUseCase.DeployConfig.PathToProject = adapter.NewFilesystemAdapter().CleanPath(postCreateDeploymentRequest.PathToSource)

		application.DeployApplication(deployUseCase)

		// h.databaseAdapter.SaveState(h.deployConfig)
		return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is deploy"})
	}
}
