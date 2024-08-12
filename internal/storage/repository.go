package storage

type Repository interface {
	GetAllMetrics() map[string]map[string]string

	GetMetric(mtype string, name string) (string, bool)

	AddGauge(name string, value float64)
	AddCounter(name string, value int64)

	SetGauges(gauges map[string]float64)
}
