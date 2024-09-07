package storage

import (
	"encoding/json"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
	rwm      sync.RWMutex
	filepath string
}

func NewMemStorage() *MemStorage {
	ms := new(MemStorage)

	return ms
}

func (ms *MemStorage) Init(restore bool, filepath string) error {
	ms.Gauges = make(map[string]float64)
	ms.Counters = make(map[string]int64)
	ms.filepath = filepath

	if restore {
		file, err := os.OpenFile(ms.filepath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		if len(data) != 0 {
			var result []models.Metric

			jsonErr := json.Unmarshal(data, &result)
			if jsonErr != nil {
				return err
			}

			for _, elem := range result {
				if elem.MType == "gauge" {
					ms.AddGauge(elem.ID, *elem.Value)
				} else {
					ms.AddCounter(elem.ID, *elem.Delta)
				}
			}
		}
	}

	return nil
}

func (ms *MemStorage) SetGauges(gauges map[string]float64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()
	for name, value := range gauges {
		ms.Gauges[name] = value
	}
}

func (ms *MemStorage) AddGauge(name string, value float64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()
	ms.Gauges[name] = value
}

func (ms *MemStorage) AddCounter(name string, value int64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()
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

func (ms *MemStorage) toMetricArray() []models.Metric {
	var metrics []models.Metric

	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	for name, elem := range ms.Gauges {
		m := models.Metric{ID: name, MType: "gauge", Value: &elem}
		metrics = append(metrics, m)
	}
	for name, elem := range ms.Counters {
		m := models.Metric{ID: name, MType: "counter", Delta: &elem}
		metrics = append(metrics, m)
	}
	return metrics
}

func (ms *MemStorage) Save(interval int) error {
	for {
		runtime.Gosched()

		data := ms.toMetricArray()
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return err
		}

		if len(data) > 0 {
			file, err := os.OpenFile(ms.filepath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = file.Write(jsonData)
			if err != nil {
				return err
			}

			logger.Log.Info("Save into file")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
