package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
	"strings"
)

type setGlobalParameterRequest struct {
	Parameter string `json:"parameter" validate:"required"`
	Value     string `json:"value" validate:"required"`
}

func setGlobalParameter(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req setGlobalParameterRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		if err := service.SetGlobalParameter(strings.TrimSpace(req.Parameter), strings.TrimSpace(req.Value)); err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to set global parameter: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	}
}
