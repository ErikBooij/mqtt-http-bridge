package server

import (
	"fmt"
	"github.com/blues/jsonata-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"text/template"
)

type validationRequest struct {
	ValidationType string `json:"type" validate:"required,oneof=jsonata template"`
	Subject        string `json:"subject" validate:"required"`
}

func validate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req validationRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

		switch req.ValidationType {
		case "jsonata":
			if err := validateJsonata(req.Subject); err != nil {
				return c.JSON(http.StatusOK, map[string]any{"error": err.Error()})
			}

			return c.JSON(http.StatusOK, nil)
		case "template":
			if err := validateTemplate(req.Subject); err != nil {
				return c.JSON(http.StatusOK, map[string]any{"error": err.Error()})
			}

			return c.JSON(http.StatusOK, nil)
		default:
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid validation type: %s", req.ValidationType))
		}
	}
}

func validateJsonata(filter string) error {
	_, err := jsonata.Compile(filter)

	if err == nil {
		return nil
	}

	return fmt.Errorf("invalid jsonata: %w", err)
}

func validateTemplate(templateString string) error {
	_, err := template.New("test").Parse(templateString)

	if err == nil {
		return nil
	}

	errMessage := err.Error()

	regex := regexp.MustCompile(`^template: test:([0-9]+):`)

	if regex.MatchString(errMessage) {
		errMessage = regex.ReplaceAllString(errMessage, "")
	}

	return fmt.Errorf("invalid template: %s", errMessage)
}
