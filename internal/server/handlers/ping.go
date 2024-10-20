package handlers

import (
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/storage"
)

type PingHandler struct {
	Ms storage.Repository
}

func (ph PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ping, ok := (ph.Ms).(storage.Ping)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := ping.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
	}
}
