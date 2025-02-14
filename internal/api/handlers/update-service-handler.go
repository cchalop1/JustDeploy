package handlers

import (
	"net/http"
	"regexp"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo/v4"
)

func isValidSubdomain(subdomain string) bool {
	regex := `^([a-zA-Z0-9-]{1,63})?$`
	re := regexp.MustCompile(regex)
	return re.MatchString(subdomain)
}

func UpdateServiceHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serviceToUpdate := domain.Service{}

		err := c.Bind(&serviceToUpdate)

		if !isValidSubdomain(serviceToUpdate.ExposeSettings.SubDomain) {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: "Invalid subdomain"})
		}

		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		service, err := application.UpdateService(deployService, serviceToUpdate)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, service)
	}
}
