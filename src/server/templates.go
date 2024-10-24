package server

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"mqtt-http-bridge/src/config"
	"mqtt-http-bridge/src/frontend"
	"net/http"
)

var templatePaths = []string{
	"*.gohtml",
	"includes/*.gohtml",
}

func newTemplateRenderer(cfg *config.Config) templateRenderer {
	if cfg.IsDevelopment() {
		return &developmentRenderer{directory: "src/frontend/templates"}
	}

	fs := frontend.Templates()
	tpl := template.New("")

	for _, path := range templatePaths {
		t, err := tpl.ParseFS(fs, path)

		if err != nil {
			panic(err)
		}

		tpl = t
	}

	return &productionRenderer{tpl: tpl}
}

type templateRenderer interface {
	Render(c echo.Context, template string, data interface{}) error
}

type developmentRenderer struct {
	directory string
}

func (r *developmentRenderer) Render(c echo.Context, templateName string, data interface{}) error {
	tpl := template.New(templateName)

	for _, path := range templatePaths {
		t, err := tpl.ParseGlob(r.directory + "/" + path)

		if err != nil {
			panic(err)
		}

		tpl = t
	}

	if err := tpl.ExecuteTemplate(c.Response().Writer, templateName, data); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}

type productionRenderer struct {
	tpl *template.Template
}

func (r *productionRenderer) Render(c echo.Context, template string, data interface{}) error {
	err := r.tpl.ExecuteTemplate(c.Response().Writer, template, data)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}
