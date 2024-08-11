package storage

import "strconv"

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (ms *MemStorage) SetGauges(gauges map[string]float64) {
	for name, value := range gauges {
		ms.Gauges[name] = value
	}
}

func (ms *MemStorage) AddGauge(name string, value float64) {
	ms.Gauges[name] = value
}

func (ms *MemStorage) AddCounter(name string, value int64) {
	ms.Counters[name] += value
}

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	val, ok := ms.Gauges[name]

	return val, ok
}

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	val, ok := ms.Counters[name]
	return val, ok
}

func (ms *MemStorage) GetAllMetrics() map[string]map[string]string {
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
