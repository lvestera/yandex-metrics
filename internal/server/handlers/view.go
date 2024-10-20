package handlers

import (
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

type ViewHandler struct {
	Ms     storage.Repository
	Format adapters.Format
}

func (mh ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m, err := mh.Format.ParseViewRequest(r)
	contentType := mh.Format.ContentType()

	w.Header().Add("Content-Type", contentType)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
		return
	}

	m, err = mh.Ms.GetMetric(m.MType, m.ID)

	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	responseBody, err := mh.Format.BuildViewResponseBody(m)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
