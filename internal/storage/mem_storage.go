package storage

import (
	"strconv"
	"sync"
)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
	rwm      sync.RWMutex
}

func NewMemStorage() *MemStorage {
	ms := new(MemStorage)
	ms.Gauges = make(map[string]float64)
	ms.Counters = make(map[string]int64)

	return ms
}

func (ms *MemStorage) SetGauges(gauges map[string]float64) {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	for name, value := range gauges {
		ms.Gauges[name] = value
	}
}

func (ms *MemStorage) AddGauge(name string, value float64) {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	ms.Gauges[name] = value
}

func (ms *MemStorage) AddCounter(name string, value int64) {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	ms.Counters[name] += value
}

func (ms *MemStorage) GetAllMetrics() map[string]map[string]string {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	result := make(map[string]map[string]string)

	result["gauge"] = make(map[string]string)
	result["counter"] = make(map[string]string)
	for name, val := range ms.Gauges {
		result["gauge"][name] = strconv.FormatFloat(val, 'f', -1, 64)
	}

	for name, val := range ms.Counters {
		result["counter"][name] = strconv.FormatInt(val, 10)
	}

	return result
}

func (ms *MemStorage) GetMetric(mtype string, name string) (string, bool) {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	switch mtype {
	case "gauge":
		val, ok := ms.Gauges[name]

		if !ok {
			return "", ok
		}
		return strconv.FormatFloat(val, 'f', -1, 64), ok
	case "counter":
		val, ok := ms.Counters[name]

		if !ok {
			return "", ok
		}

		return strconv.FormatInt(val, 10), ok
	}

	return "", false
}
