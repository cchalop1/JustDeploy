package handlers

import (
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func SubscriptionCreateServer(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := c.Response()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for e := range deployService.EventAdapter.EventServer {
			e.MarshalToSseEvent(w)
			w.Flush()
		}

		return nil
	}
}
