package server

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"reflect"
	"strings"
)

type customValidator struct {
	validator *validator.Validate
}

func newValidator() echo.Validator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &customValidator{validator: validate}
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}

	return nil
}
