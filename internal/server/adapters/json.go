package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/models"
)

type Json struct{}

func (f Json) ParseUpdateRequest(r *http.Request) (models.Metric, error) {
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

func (f Json) ParseViewRequest(r *http.Request) (models.Metric, error) {
	var m models.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(buf.Bytes(), &m)

	return m, err
}

func (f Json) BuildUpdateResponseBody(m models.Metric) ([]byte, error) {
	return json.Marshal(m)
}

func (f Json) BuildViewResponseBody(m models.Metric) ([]byte, error) {
	return json.Marshal(m)
}

func (f Json) ContentType() string {
	return "application/json"
}
