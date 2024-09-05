package adapters

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/models"
)

type Http struct{}

func (f Http) ParseUpdateRequest(r *http.Request) (models.Metric, error) {
	m := models.Metric{ID: chi.URLParam(r, "name"), MType: chi.URLParam(r, "mtype")}

	err := m.SetValue(chi.URLParam(r, "value"))

	return m, err
}

func (f Http) ParseViewRequest(r *http.Request) (models.Metric, error) {

	m := models.Metric{ID: chi.URLParam(r, "name"), MType: chi.URLParam(r, "mtype")}

	return m, nil
}

func (f Http) BuildUpdateResponseBody(_ models.Metric) ([]byte, error) {
	return []byte(nil), nil
}
func (f Http) BuildViewResponseBody(m models.Metric) ([]byte, error) {

	value, err := m.GetValue()
	return ([]byte)(value), err
}

func (f Http) ContentType() string {
	return "text/plain"
}
