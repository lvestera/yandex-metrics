package main

import (
	"github.com/lvestera/yandex-metrics/internal/agent"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

func main() {

	parseFlags()
	metric := &storage.MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}

	c := &agent.MetricClient{
		Host: addr,
	}

	go agent.Update(metric, pollInterval)
	agent.Send(metric, c, reportInterval)
}
