package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type updateSubscriptionRequest struct {
	Name  string `json:"name" validate:"required_without=SubscriptionTemplateID"`
	Topic string `json:"topic" validate:"required_without=SubscriptionTemplateID"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	Method  string            `json:"method" validate:"required_without=SubscriptionTemplateID,omitempty,oneof=GET POST PUT PATCH DELETE"`
	URL     string            `json:"url" validate:"required_without=SubscriptionTemplateID"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func updateSubscription(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req updateSubscriptionRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		sub, err := service.UpdateSubscription(subscription.Subscription{
			ID: c.Param("id"),

			Name:  req.Name,
			Topic: req.Topic,

			Extract: req.Extract,
			Filter:  req.Filter,

			Method:  req.Method,
			URL:     req.URL,
			Headers: req.Headers,
			Body:    req.Body,
		})

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to update subscription: %w", err))
		}

		return c.JSON(http.StatusCreated, map[string]any{"subscription": subscriptionToResponse(sub)})
	}
}
