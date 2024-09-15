package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

func listGlobalParameters(service subscription.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		params, err := service.GetGlobalParameters()

		if err != nil {
			return ErrorResponse(c, mapErrorCode(err), fmt.Errorf("failed to list global parameters: %w", err))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"parameters": params})
	}
}
