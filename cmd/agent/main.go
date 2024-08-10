package main

import (
	"github.com/lvestera/yandex-metrics/internal/agent"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

func main() {

	metric := &storage.MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}

	c := &agent.MetricClient{
		Host: "http://localhost:8080",
	}

	go agent.Update(metric)
	agent.Send(metric, c)
}
