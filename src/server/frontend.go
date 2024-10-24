package server

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/datastore"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

func redirect(path string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, path)
	}
}

func globalParameters(tpl templateRenderer) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := sharedGlobalTemplateParameters()

		params.JSFiles = []string{"/assets/global-parameters.js"}
		params.PageTitle = "Global Parameters"

		return tpl.Render(c, "global-parameters.gohtml", params.WithCurrentMenuItem("Global Parameters"))
	}
}

func globalParameterCreate(tpl templateRenderer) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := sharedGlobalTemplateParameters()

		params.JSFiles = []string{"/assets/global-parameter.js"}
		params.PageTitle = "Global Parameter"

		return tpl.Render(c, "global-parameter.gohtml", params.WithCurrentMenuItem("Global Parameters"))
	}
}

func globalParameterUpdate(tpl templateRenderer, service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		globalParameters, err := service.GetGlobalParameters()

		if err != nil {
			return errorPage(c, tpl, 500, "Internal server error")
		}

		if _, ok := globalParameters[key]; !ok {
			return errorPage(c, tpl, 404, "Global parameter not found")
		}

		params := sharedGlobalTemplateParameters()

		params.JSFiles = []string{"/assets/global-parameter.js"}
		params.PageTitle = "Global Key"
		params.Data = map[string]any{
			"GlobalParameterKey": key,
		}

		return tpl.Render(c, "global-parameter.gohtml", params.WithCurrentMenuItem("Global Parameters"))
	}
}

func subscriptions(tpl templateRenderer) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := sharedGlobalTemplateParameters()

		params.JSFiles = []string{"/assets/subscriptions.js"}
		params.PageTitle = "Subscriptions"

		return tpl.Render(c, "subscriptions.gohtml", params.WithCurrentMenuItem("Subscriptions"))
	}
}

func subscriptionCreate(tpl templateRenderer) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := sharedGlobalTemplateParameters()

		params.JSFiles = []string{"/assets/subscription.js"}
		params.PageTitle = "Subscription"

		return tpl.Render(c, "subscription.gohtml", params.WithCurrentMenuItem("Subscriptions"))
	}
}

func subscriptionUpdate(tpl templateRenderer, service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		subID := c.Param("id")

		_, err := service.GetSubscription(subID)

		switch {
		case errors.Is(err, datastore.ErrSubscriptionNotFound):
			return errorPage(c, tpl, 404, "Subscription not found")
		case err != nil:
			return errorPage(c, tpl, 500, "Internal server error")
		}

		params := sharedGlobalTemplateParameters()

		params.JSFiles = []string{"/assets/subscription.js"}
		params.PageTitle = "Subscription"
		params.Data = map[string]any{
			"SubscriptionID": subID,
		}

		return tpl.Render(c, "subscription.gohtml", params.WithCurrentMenuItem("Subscriptions"))
	}
}

func errorPage(c echo.Context, tpl templateRenderer, httpStatusCode int, error string) error {
	params := sharedGlobalTemplateParameters()

	params.PageTitle = fmt.Sprintf("%d - Error", httpStatusCode)
	params.Data = error

	return tpl.Render(c, "error.gohtml", params)
}

type globalTemplateParameters struct {
	PageTitle string
	JSFiles   []string
	MenuItems []MenuItem
	Data      any
}

type MenuItem struct {
	Title   string
	URL     string
	Current bool
}

func sharedGlobalTemplateParameters() globalTemplateParameters {
	return globalTemplateParameters{
		JSFiles: make([]string, 0),
		MenuItems: []MenuItem{
			{"Subscriptions", "/subscriptions", false},
			{"Templates", "/templates", false},
			{"Global Parameters", "/global-parameters", false},
			{"MQTT", "/mqtt", false},
		},
	}
}

func (gtp globalTemplateParameters) WithCurrentMenuItem(title string) globalTemplateParameters {
	for i, item := range gtp.MenuItems {
		if item.Title == title {
			gtp.MenuItems[i].Current = true
		}
	}

	return gtp
}
