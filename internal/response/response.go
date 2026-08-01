package response

import (
	"encoding/json"
	"net/http"
)

func WithStatusOnly(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func WithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&data)
}

func WithError(w http.ResponseWriter, err error, code int) {
	responseError := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: err.Error(),
	}
	WithJSON(w, responseError, statusCode)
}

