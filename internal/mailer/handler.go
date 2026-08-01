package mailer

import (
	"net/http"
)

type (
	Handler struct {

	}
)

func (h Handler) RegisterToken(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) DeleteTokens(w http.ResponseWriter, r *http.Request) {
}

