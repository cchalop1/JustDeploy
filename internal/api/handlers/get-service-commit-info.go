package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Get service commit information
// @Description Retrieves the last deployed commit information for a GitHub service
// @Tags service
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Success 200 {object} application.ServiceCommitInfo "Commit information"
// @Success 204 "No commit information available (not a GitHub service)"
// @Failure 400 {object} dto.ResponseApi "Bad request"
// @Failure 404 {object} dto.ResponseApi "Service not found"
// @Failure 500 {object} dto.ResponseApi "Internal server error"
// @Router /api/v1/service/{serviceId}/commit-info [get]
func GetServiceCommitInfo(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serviceId := c.Param("serviceId")
		if serviceId == "" {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: "Service ID is required"})
		}

		commitInfo, err := application.GetServiceCommitInfo(deployService, serviceId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		// Si commitInfo est nil, cela signifie que ce n'est pas un service GitHub
		if commitInfo == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return c.JSON(http.StatusOK, commitInfo)
	}
}
