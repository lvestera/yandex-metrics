package handlers

import (
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

type UpdateHandler struct {
	Ms     storage.Repository
	Format adapters.Format
}

func (uh UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, err := uh.Format.ParseUpdateRequest(r)
	contentType := uh.Format.ContentType()

	w.Header().Add("Content-Type", contentType)

	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = uh.Ms.AddMetric(m)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	m, err = uh.Ms.GetMetric(m.MType, m.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// value, ok := uh.Ms.GetMetric(m.MType, m.ID)
	// if !ok {
	// 	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	// 	return
	// }
	// m.SetValue(value)

	responseBody, err := uh.Format.BuildUpdateResponseBody(m)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
