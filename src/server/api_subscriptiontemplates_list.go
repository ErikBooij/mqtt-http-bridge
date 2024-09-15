package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"mqtt-http-bridge/src/utilities"
	"net/http"
)

func listSubscriptionTemplates(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		subs, err := service.GetSubscriptionTemplates()

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to list subscription templates: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"subscriptions": utilities.MapSlice(subs, subscriptionTemplateToResponse)})
	}
}
