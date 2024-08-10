package main

import (
	. "github.com/lvestera/yandex-metrics/internal/agent"
	. "github.com/lvestera/yandex-metrics/internal/storage"
)

func main() {

	metric := &MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}

	c := &MetricClient{
		Host: "http://localhost:8080",
	}

	go Update(metric)
	Send(metric, c)
}
