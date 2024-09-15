package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type deleteSubscriptionRequest struct {
	ID string `param:"id" validate:"required"`
}

func deleteSubscription(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req deleteSubscriptionRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		if err := service.DeleteSubscription(req.ID); err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to delete subscription: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]any{"status": "success"})
	}
}
