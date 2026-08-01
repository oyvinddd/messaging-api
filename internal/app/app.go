package app

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
)

const (
	/* token ruotes */
	registerTokenRoute 	= "POST /api/v1/push/register"
	deleteTokensRoute 	= "DELETE /api/v1/push/delete"
	/* email routes */
	sendEmailRoute 		= "POST /api/v1/email/send"
)

type (
	App struct {
		addr	string
		server *http.Server
		mux    *http.ServeMux
	}
)

func New() *App {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load environment: %v\n", err)
	}

	addr := fmt.Sprintf(":%s", os.Getenv("LISTENING_PORT"))
	mux := http.NewServeMux()

	app := &App{
		addr: addr,
		mux: mux,
		server:	&http.Server{
			Handler: mux,
			Addr: addr,
		},
	}

	return app
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

