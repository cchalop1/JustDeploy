package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func PostAddDomainToServerById(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		newDeployDto := dto.NewDomain{}
		err := c.Bind(&newDeployDto)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		serverId := c.Param("id")

		err = application.AddDomainToServer(deployService, newDeployDto, serverId)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, true)
	}
}
