package storage

import "github.com/lvestera/yandex-metrics/internal/models"

type Repository interface {
	GetMetrics() ([]models.Metric, error)
	GetMetric(mtype string, name string) (models.Metric, error)

	AddMetrics(metrics []models.Metric) (int, error)
	AddMetric(m models.Metric) (bool, error)

	SetGauges(gauges map[string]float64)

	Save(interval int) error
}
