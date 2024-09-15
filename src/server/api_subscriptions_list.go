package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"mqtt-http-bridge/src/utilities"
	"net/http"
)

func listSubscriptions(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		subs, err := service.GetSubscriptions()

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to list subscriptions: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"subscriptions": utilities.MapSlice(subs, subscriptionToResponse)})
	}
}
