package server

import (
	"errors"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

func mapErrorCode(err error) int {
	switch {
	case errors.Is(err, subscription.ErrMissingRequiredParametersForTemplate):
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
