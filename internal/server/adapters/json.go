package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/models"
)

type JSON struct{}

func (f JSON) ParseUpdateRequest(r *http.Request) (models.Metric, error) {
	var m models.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(buf.Bytes(), &m)

	if m.ID == "" || m.MType == "" {
		return m, errors.New("some fields missing")
	}

	return m, err
}

func (f JSON) ParseUpdateBatchRequest(r *http.Request) ([]models.Metric, error) {
	var metrics []models.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return metrics, err
	}

	err = json.Unmarshal(buf.Bytes(), &metrics)

	return metrics, err
}

func (f JSON) ParseViewRequest(r *http.Request) (models.Metric, error) {
	var m models.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(buf.Bytes(), &m)

	return m, err
}

func (f JSON) BuildUpdateResponseBody(m models.Metric) ([]byte, error) {
	return json.Marshal(m)
}

func (f JSON) BuildViewResponseBody(m models.Metric) ([]byte, error) {
	return json.Marshal(m)
}

func (f JSON) ContentType() string {
	return "application/json"
}
