package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func page(template string) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := tpl.ExecuteTemplate(c.Response().Writer, template, nil)

		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}
