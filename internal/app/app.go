package app

import (
	"net/http"
)

type (
	App struct {
		server http.Server
	}
)

func New() *App {
	return &App{}
}

func (app *App) Run() error {
	return app.server.ListenAndServe()
}

