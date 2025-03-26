package agent

import (
	"math/rand/v2"
	"reflect"
	"runtime"

	"github.com/lvestera/yandex-metrics/internal/models"
)

var MetricsName = [...]string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups",
	"MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
	"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys",
	"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

func collectMetrics(rtm *runtime.MemStats) []models.Metric {
	runtime.ReadMemStats(rtm)

	metrics := make([]models.Metric, 0)

	var m models.Metric
	var val float64
	for _, mname := range MetricsName {
		r := reflect.ValueOf(*rtm)
		f := reflect.Indirect(r).FieldByName(mname)
		if f.CanUint() {
			val = float64(f.Uint())
		}
		if f.CanFloat() {
			val = float64(f.Float())
		}

		m = models.Metric{
			ID:    mname,
			MType: "gauge",
			Value: &val,
		}

		metrics = append(metrics, m)
	}

	val = rand.Float64()
	m = models.Metric{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &val,
	}

	metrics = append(metrics, m)

	return metrics
}
