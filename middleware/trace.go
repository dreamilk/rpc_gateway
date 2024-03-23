package middleware

import (
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

func (h handler) WithRequestID() handler {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
