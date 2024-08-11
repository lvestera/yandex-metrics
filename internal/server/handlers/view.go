package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

type ViewHandler struct {
	Ms storage.Repository
}

func (mh ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mtype")
	mName := chi.URLParam(r, "name")

	value, ok := mh.Ms.GetMetric(mType, mName)

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	io.WriteString(w, value)
}
