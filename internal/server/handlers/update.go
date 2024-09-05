package handlers

import (
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/adapters"
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
		http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "gauge":
		uh.updateGauge(m)
	case "counter":
		uh.updateCounter(m)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	value, ok := uh.Ms.GetMetric(m.MType, m.ID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	m.SetValue(value)

	responseBody, err := uh.Format.BuildUpdateResponseBody(m)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func (uh UpdateHandler) updateGauge(m models.Metric) {
	uh.Ms.AddGauge(m.ID, *m.Value)
}

func (uh UpdateHandler) updateCounter(m models.Metric) {
	uh.Ms.AddCounter(m.ID, *m.Delta)
}
