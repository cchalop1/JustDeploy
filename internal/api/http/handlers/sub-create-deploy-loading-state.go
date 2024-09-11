package handlers

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func SubscriptionCreateDeployLoadingState(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := c.Response()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		for {
			select {
			case <-c.Request().Context().Done():
				return nil
			case event := <-deployService.EventAdapter.EventDeployWrapper:
				fmt.Println("Sending event", event)
				event.MarshalToSseEvent(w)
				w.Flush()
			}
		}
	}
}
