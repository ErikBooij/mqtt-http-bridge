package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type addSubscriptionRequest struct {
	Name  string `json:"name" validate:"required"`
	Topic string `json:"topic" validate:"required"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	Method  string            `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	URL     string            `json:"url" validate:"required"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func addSubscription(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req addSubscriptionRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		}

		sub, err := service.AddSubscription(subscription.Subscription{
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
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to add subscription: %w", err))
		}

		return c.JSON(http.StatusCreated, map[string]any{"subscription": subscriptionToResponse(sub)})
	}
}
