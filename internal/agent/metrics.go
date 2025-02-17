package agent

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/shirou/gopsutil/v4/mem"
)

var MetricsName = [...]string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups",
	"MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
	"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys",
	"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

func UpdateMainMetrics(m storage.Repository, interval int) {

	var rtm runtime.MemStats
	for {
		runtime.ReadMemStats(&rtm)

		m.AddMetrics(collectMetrics(&rtm))

		var pollMetricValue int64 = 1

		m.AddMetric(models.Metric{
			ID:    "PollCount",
			MType: "counter",
			Delta: &pollMetricValue,
		})

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func UpdateAdditionalMetrics(m storage.Repository, interval int) {

	var val float64
	for {
		v, _ := mem.VirtualMemory()
		metrics := make([]models.Metric, 0)

		val = float64(v.Total)
		metrics = append(metrics, models.Metric{
			ID:    "TotalMemory",
			MType: "gauge",
			Value: &val,
		})

		val = float64(v.Free)
		metrics = append(metrics, models.Metric{
			ID:    "FreeMemory",
			MType: "gauge",
			Value: &val,
		})
		val = float64(v.UsedPercent)
		metrics = append(metrics, models.Metric{
			ID:    "CPUutilization1",
			MType: "gauge",
			Value: &val,
		})

		m.AddMetrics(metrics)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// func Send(m storage.Repository, c MClient, interval int, key string, rateLimit int) {

// 	for {
// 		metrics, err := m.GetMetrics()
// 		if err != nil {
// 			logger.Log.Info("Get metrics failed")
// 		}

// 		err = c.SendBatchUpdate(metrics, key)
// 		if err != nil {
// 			logger.Log.Info(fmt.Sprint("Sending the batch of ", len(metrics), "metrics failed: ", err.Error()))
// 		}

//			time.Sleep(time.Duration(interval) * time.Second)
//		}
//	}
func Send(jobs <-chan struct{}, m storage.Repository, c MClient, key string) {

	for range jobs {
		metrics, err := m.GetMetrics()
		if err != nil {
			logger.Log.Info("Get metrics failed")
		}

		err = c.SendBatchUpdate(metrics, key)
		if err != nil {
			logger.Log.Info(fmt.Sprint("Sending the batch of ", len(metrics), "metrics failed: ", err.Error()))
		}
	}
}
func PrepareSend(jobs chan<- struct{}, interval int) {
	for {
		jobs <- struct{}{}
		time.Sleep(time.Duration(interval) * time.Second)
	}
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
