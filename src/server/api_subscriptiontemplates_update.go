package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type updateSubscriptionTemplateRequest struct {
	ID string `param:"id" validate:"required"`

	Name  string `json:"name" validate:"required"`
	Topic string `json:"topic" validate:"required"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	Method       string            `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	URL          string            `json:"url" validate:"required"`
	Headers      map[string]string `json:"headers"`
	BodyTemplate string            `json:"bodyTemplate"`

	RequiredParameters []string `json:"requiredParameters" validate:"required"`
}

func updateSubscriptionTemplate(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req updateSubscriptionTemplateRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		subTemp, err := service.UpdateSubscriptionTemplate(subscription.SubscriptionTemplate{
			Subscription: subscription.Subscription{
				ID: req.ID,

				Name:  req.Name,
				Topic: req.Topic,

				Extract: req.Extract,
				Filter:  req.Filter,

				Method:       req.Method,
				URL:          req.URL,
				Headers:      req.Headers,
				BodyTemplate: req.BodyTemplate,
			},
			RequiredParameters: req.RequiredParameters,
		})

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to update subscription template: %w", err))
		}

		return c.JSON(http.StatusCreated, map[string]any{"subscriptionTemplate": subscriptionTemplateToResponse(subTemp)})
	}
}
