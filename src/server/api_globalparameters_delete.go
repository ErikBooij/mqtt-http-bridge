package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
	"strings"
)

type deleteGlobalParameterRequest struct {
	Parameter string `param:"parameter" validate:"required"`
}

func deleteGlobalParameter(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req deleteGlobalParameterRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		if err := service.DeleteGlobalParameter(strings.TrimSpace(req.Parameter)); err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to set global parameter: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	}
}
