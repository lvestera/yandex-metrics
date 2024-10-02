package handlers

import "net/http"

type PingHandler struct{}

func (mh PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
