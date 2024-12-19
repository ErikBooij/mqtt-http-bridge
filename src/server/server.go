package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"mqtt-http-bridge/src/config"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type HTTPServer interface {
	Start(address string) error
}

func New(service subscription.Service, cfg *config.Config) HTTPServer {
	server := echo.New()
	server.Binder = newBinder()
	server.Validator = newValidator()

	logger := log.New("")
	logger.SetOutput(io.Discard)
	server.Logger = logger

	server.GET("/assets/*", assets(cfg))

	server.Any("/*", app())

	api := server.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})
	api.POST("/validate", validate())

	api.DELETE("/subscriptions/:id", deleteSubscription(service))
	api.GET("/subscriptions/:id", getSubscription(service))
	api.PUT("/subscriptions/:id", updateSubscription(service))
	api.GET("/subscriptions", listSubscriptions(service))
	api.POST("/subscriptions", addSubscription(service))

	api.DELETE("/global-parameters/:parameter", deleteGlobalParameter(service))
	api.GET("/global-parameters", listGlobalParameters(service))
	api.POST("/global-parameters", setGlobalParameter(service))

	api.Any("/*", apiError(http.StatusNotFound, "Not Found"))

	return server
}

func apiError(code int, message string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(code, map[string]any{"error": message})
	}
}
