package agent

import (
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

	. "github.com/lvestera/yandex-metrics/internal/storage"
)

var Metrics_name = [...]string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups",
	"MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
	"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys",
	"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

func Update(m Repository) {

	var rtm runtime.MemStats
	for {
		runtime.Gosched()
		runtime.ReadMemStats(&rtm)

		m.SetGauges(collectMetrics(rtm))
		m.AddCounter("PollCount", 1)

		time.Sleep(2 * time.Second)
	}
}

func Send(m Repository, c MClient) {

	for {
		runtime.Gosched()

		for mtype, row := range m.GetAllMetrics() {
			for name, svalue := range row {
				err := c.SendUpdate(mtype, name, svalue)
				if err != nil {
					return
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func collectMetrics(rtm runtime.MemStats) map[string]float64 {
	runtime.ReadMemStats(&rtm)

	metrics_value := make(map[string]float64)

	for _, mname := range Metrics_name {
		r := reflect.ValueOf(rtm)
		f := reflect.Indirect(r).FieldByName(mname)
		if f.CanUint() {
			metrics_value[mname] = float64(f.Uint())
		}
		if f.CanFloat() {
			metrics_value[mname] = float64(f.Float())
		}

	}

	metrics_value["RandomValue"] = rand.Float64()

	return metrics_value

}
