package handlers

// import (
// 	"net/http"

// 	"cchalop1.com/deploy/internal/api/service"
// 	"cchalop1.com/deploy/internal/application"
// 	"github.com/labstack/echo/v4"
// )

// type CreateDatabaseRequest struct {
// 	DatabaseName string `json:"databaseName"`
// }

// func PostCreateDatabaseHandler(deployService *service.DeployService) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var req CreateDatabaseRequest
// 		if err := c.Bind(&req); err != nil {
// 			return c.JSON(http.StatusBadRequest, map[string]string{
// 				"error": "Invalid request format",
// 			})
// 		}

// 		// Validate required fields
// 		if req.DatabaseName == "" {
// 			return c.JSON(http.StatusBadRequest, map[string]string{
// 				"error": "DatabaseName is required",
// 			})
// 		}

// 		service, err := application.CreateServiceFromDatabase(deployService, req.DatabaseName)
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, map[string]string{
// 				"error": err.Error(),
// 			})
// 		}

// 		return c.JSON(http.StatusCreated, service)
// 	}
// }
