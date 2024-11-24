package server

import (
	"fmt"
	"github.com/blues/jsonata-go"
	"github.com/labstack/echo/v4"
	"net/http"
)

type validationRequest struct {
	ValidationType string `json:"type" validate:"required,oneof=extract filter"`
	Subject        string `json:"subject" validate:"required"`
}

func validate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req validationRequest

		if err := c.Bind(&req); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

		switch req.ValidationType {
		case "extract":
			if err := validateJsonata(req.Subject); err != nil {
				return c.JSON(http.StatusOK, map[string]any{"error": err.Error()})
			}

			return c.JSON(http.StatusOK, nil)
		case "filter":
			if err := validateJsonata(req.Subject); err != nil {
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
