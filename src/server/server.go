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
	tplRenderer := newTemplateRenderer(cfg)

	server := echo.New()
	server.Binder = newBinder()
	server.Validator = newValidator()

	logger := log.New("")
	logger.SetOutput(io.Discard)
	server.Logger = logger

	server.GET("/assets/*", assets(cfg))

	server.GET("/", redirect("/subscriptions"))

	server.GET("/subscriptions", subscriptions(tplRenderer))
	server.GET("/subscriptions/:id", subscriptionEdit(tplRenderer, service))
	server.GET("/new-subscription", subscriptionCreate(tplRenderer, service))

	server.GET("/global-parameters", globalParameters(tplRenderer))
	server.GET("/global-parameters/:key", globalParameterEdit(tplRenderer, service))
	server.GET("/new-global-parameter", globalParameterCreate(tplRenderer))

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

	return server
}
