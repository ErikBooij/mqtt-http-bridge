package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var assets embed.FS

func Assets() fs.FS {
	return SubOrPanic(assets, "dist")
}

//go:embed templates
var templates embed.FS

func Templates() fs.FS {
	return SubOrPanic(templates, "templates")
}

func SubOrPanic(f embed.FS, name string) fs.FS {
	fs, err := fs.Sub(f, name)

	if err != nil {
		panic(err)
	}

	return fs
}
