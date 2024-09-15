package server

import (
	"github.com/labstack/echo/v4"
)

type customBinder struct {
	defaultBinder echo.Binder
}

func newBinder() echo.Binder {
	return &customBinder{
		defaultBinder: new(echo.DefaultBinder),
	}
}

func (b customBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.defaultBinder.Bind(i, c); err != nil {
		return err
	}

	return c.Validate(i)
}
