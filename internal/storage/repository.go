package storage

import (
	"context"

	"github.com/lvestera/yandex-metrics/internal/models"
)

type Repository interface {
	GetMetrics(ctx context.Context) ([]models.Metric, error)
	GetMetric(ctx context.Context, mtype string, name string) (models.Metric, error)

	AddMetrics(metrics []models.Metric) (int, error)
	AddMetric(m models.Metric) error

	Save(ctx context.Context, interval int) error
}
