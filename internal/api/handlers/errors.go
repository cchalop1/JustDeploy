package handlers

import (
	"cchalop1.com/deploy/internal/domain"
	"errors"
	"github.com/labstack/echo/v4"
)

func operationError(err error) error {
	if errors.Is(err, domain.ErrConflict) || errors.Is(err, domain.ErrNotFound) {
		return err
	}
	return echo.NewHTTPError(400, err.Error())
}
