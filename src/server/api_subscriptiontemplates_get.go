package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type getSubscriptionTemplateRequest struct {
	ID string `param:"id" validate:"required"`
}

func getSubscriptionTemplate(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req getSubscriptionTemplateRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		sub, err := service.GetSubscriptionTemplate(req.ID)

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to get subscription: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"subscriptionTemplate": subscriptionTemplateToResponse(sub)})
	}
}
