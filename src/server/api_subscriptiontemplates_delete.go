package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type deleteSubscriptionTemplateRequest struct {
	ID string `param:"id" validate:"required"`
}

func deleteSubscriptionTemplate(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req deleteSubscriptionTemplateRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		if err := service.DeleteSubscriptionTemplate(req.ID); err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to delete subscription template: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]any{"status": "success"})
	}
}
