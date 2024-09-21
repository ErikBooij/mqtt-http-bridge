package server

import (
	"html/template"
	"mqtt-http-bridge/src/frontend"
)

var tpl *template.Template

func init() {
	t, err := template.ParseFS(frontend.Templates(), "*.gohtml")

	if err != nil {
		panic(err)
	}

	tpl = t
}
