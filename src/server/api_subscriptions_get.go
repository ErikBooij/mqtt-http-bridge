package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type getSubscriptionRequest struct {
	ID string `param:"id" validate:"required"`
}

func getSubscription(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req getSubscriptionRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		sub, err := service.GetSubscription(req.ID)

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to get subscription: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"subscription": subscriptionToResponse(sub)})
	}
}
