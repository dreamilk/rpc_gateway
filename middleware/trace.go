package middleware

import (
	"net/http"

	"github.com/dreamilk/rpc_gateway/utils"
)

type handler func(http.ResponseWriter, *http.Request)

func (h handler) WithRequestID() handler {
	return func(w http.ResponseWriter, req *http.Request) {
		requestId := utils.UUID()
		req.Header.Set("request-id", requestId)
		h(w, req)
		w.Header().Add("request-id", requestId)
	}
}
