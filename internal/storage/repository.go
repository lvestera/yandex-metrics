package storage

type Repository interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)

	SetGauges(gauges map[string]float64)

	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)

	GetAllMetrics() map[string]map[string]string
}
