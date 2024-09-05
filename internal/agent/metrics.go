package agent

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

var MetricsName = [...]string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups",
	"MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
	"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys",
	"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

func Update(m storage.Repository, interval int) {

	var rtm runtime.MemStats
	for {
		runtime.Gosched()
		runtime.ReadMemStats(&rtm)

		m.SetGauges(collectMetrics(&rtm))
		m.AddCounter("PollCount", 1)

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func Send(m storage.Repository, c MClient, interval int) {

	for {
		runtime.Gosched()

		for mtype, row := range m.GetAllMetrics() {
			for name, svalue := range row {
				err := c.SendUpdate(mtype, name, svalue)
				if err != nil {
					logger.Log.Info(fmt.Sprint("Sending the", mtype, "metric", name, "failed"))
				}
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func collectMetrics(rtm *runtime.MemStats) map[string]float64 {
	runtime.ReadMemStats(rtm)

	metricsValue := make(map[string]float64)

	for _, mname := range MetricsName {
		r := reflect.ValueOf(*rtm)
		f := reflect.Indirect(r).FieldByName(mname)
		if f.CanUint() {
			metricsValue[mname] = float64(f.Uint())
		}
		if f.CanFloat() {
			metricsValue[mname] = float64(f.Float())
		}
	}

	metricsValue["RandomValue"] = rand.Float64()

	return metricsValue
}
