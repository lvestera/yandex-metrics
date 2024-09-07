package storage

import "github.com/lvestera/yandex-metrics/internal/models"

type Repository interface {
	Init(restore bool, filepath string) error

	GetMetrics() []models.Metric

	GetMetric(mtype string, name string) (string, bool)

	AddGauge(name string, value float64)
	AddCounter(name string, value int64)

	SetGauges(gauges map[string]float64)
}
