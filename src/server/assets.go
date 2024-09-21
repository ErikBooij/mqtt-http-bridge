package server

import (
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/frontend"
	"net/http"
)

func assets() echo.HandlerFunc {
	return echo.WrapHandler(http.StripPrefix("/assets", http.FileServerFS(frontend.Assets())))
}
