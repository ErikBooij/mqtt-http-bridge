package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type HTTPServer interface {
	Start(address string) error
}

func New(service subscription.Service) HTTPServer {
	server := echo.New()
	server.Binder = newBinder()
	server.Validator = newValidator()

	logger := log.New("")
	logger.SetOutput(io.Discard)
	server.Logger = logger

	api := server.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	api.DELETE("/subscriptions/:id", deleteSubscription(service))
	api.GET("/subscriptions/:id", getSubscription(service))
	api.PUT("/subscriptions/:id", updateSubscription(service))
	api.GET("/subscriptions", listSubscriptions(service))
	api.POST("/subscriptions", addSubscription(service))

	api.DELETE("/subscription-templates/:id", deleteSubscriptionTemplate(service))
	api.GET("/subscription-templates/:id", getSubscriptionTemplate(service))
	api.PUT("/subscription-templates/:id", updateSubscriptionTemplate(service))
	api.GET("/subscription-templates", listSubscriptionTemplates(service))
	api.POST("/subscription-templates", addSubscriptionTemplate(service))

	api.DELETE("/global-parameters/:parameter", deleteGlobalParameter(service))
	api.GET("/global-parameters", listGlobalParameters(service))
	api.POST("/global-parameters", setGlobalParameter(service))

	// TODO: Detach subscription from template endpoint.

	return server
}
