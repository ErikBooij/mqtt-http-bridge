package server

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ErrorResponse(c echo.Context, code int, message any) error {
	var err []string

	switch m := message.(type) {
	case error:
		err = mapError(m)
	case []error:
		switch len(m) {
		case 0:
			err = nil
		case 1:
			err = mapError(m[0])
		default:
			errs := make([]string, 0, len(m))

			for _, err := range m {
				errs = append(errs, mapError(err)...)
			}

			err = errs
		}
	case string:
		err = []string{m}
	case []string:
		err = m
	default:
		err = []string{fmt.Sprintf("%v", m)}
	}

	var rErr any

	if len(err) == 1 {
		rErr = err[0]
	} else {
		rErr = err
	}

	return c.JSON(code, map[string]interface{}{"error": rErr})
}

func mapError(err error) []string {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		errs := make([]string, 0, len(validationErrors))

		for _, validationErr := range validationErrors {
			switch validationErr.Tag() {
			case "email":
				errs = append(errs, fmt.Sprintf("Field '%s' must be a valid email address", validationErr.Field()))
			case "len":
				errs = append(errs, fmt.Sprintf("Field '%s' must be exactly %v characters long", validationErr.Field(), validationErr.Param()))
			case "oneof":
				errs = append(errs, fmt.Sprintf("Field '%s' must be one of %v", validationErr.Field(), validationErr.Param()))
			case "required", "required_without", "required_with":
				errs = append(errs, fmt.Sprintf("Field '%s' cannot be blank", validationErr.Field()))
			default:
				errs = append(errs, fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", validationErr.Field(), validationErr.Value(), validationErr.Tag(), validationErr.Param()))
			}
		}

		return errs
	}

	return []string{err.Error()}
}
