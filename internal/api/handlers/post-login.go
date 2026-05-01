package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func PostLoginHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		loginDto := dto.LoginDto{}
		if err := c.Bind(&loginDto); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{
				Message: "Invalid request format",
			})
		}

		token, err := application.Login(deployService, loginDto)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, dto.ResponseApi{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, dto.AuthResponseDto{Token: token})
	}
}
