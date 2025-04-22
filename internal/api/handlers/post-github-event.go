package handlers

import (
	"fmt"
	"io"
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func PostGithubEvent(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Read the request body
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to read request body",
			})
		}

		// Print the request body
		fmt.Println("GitHub Event Request Body:")
		fmt.Println(string(body))

		return c.JSON(http.StatusOK, true)
	}
}
