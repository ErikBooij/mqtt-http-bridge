package server

import (
	"github.com/labstack/echo/v4"
	"io/fs"
	"mqtt-http-bridge/src/config"
	"mqtt-http-bridge/src/frontend"
	"net/http"
	"os"
)

func assets(cfg *config.Config) echo.HandlerFunc {
	var dir fs.FS

	if cfg.IsDevelopment() {
		dir = os.DirFS("src/frontend/dist")
	} else {
		dir = frontend.Assets()
	}

	return echo.WrapHandler(http.StripPrefix("/assets", http.FileServer(http.FS(dir))))
}
