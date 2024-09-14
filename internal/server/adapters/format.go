package adapters

import (
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/models"
)

type Format interface {
	ParseUpdateRequest(r *http.Request) (models.Metric, error)
	ParseViewRequest(r *http.Request) (models.Metric, error)
	BuildUpdateResponseBody(m models.Metric) ([]byte, error)
	BuildViewResponseBody(m models.Metric) ([]byte, error)
	ContentType() string
}
